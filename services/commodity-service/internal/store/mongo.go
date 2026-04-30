package store

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

// CollectionName is the MongoDB collection that holds commodities.
const CollectionName = "commodities"

// Mongo is a MongoDB-backed Store. Construct it via NewMongo, which also
// ensures the unique index on `name` exists.
type Mongo struct {
	coll *mgo.Collection
	now  func() time.Time
}

// NewMongo returns a Store backed by db.Collection("commodities"). It creates
// a unique case-insensitive index on `name` so the simulation can rely on
// "one canonical record per commodity."
func NewMongo(ctx context.Context, db *mgo.Database) (*Mongo, error) {
	coll := db.Collection(CollectionName)
	model := mgo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true).SetCollation(&options.Collation{Locale: "en", Strength: 2}),
	}
	if _, err := coll.Indexes().CreateOne(ctx, model); err != nil {
		return nil, err
	}
	return &Mongo{coll: coll, now: time.Now}, nil
}

func (m *Mongo) Create(ctx context.Context, c commodity.Commodity) (commodity.Commodity, error) {
	if err := c.Validate(); err != nil {
		return commodity.Commodity{}, err
	}
	if c.ID.IsZero() {
		c.ID = commodity.NewID()
	}
	now := m.now()
	c.CreatedAt = now
	c.UpdatedAt = now
	if _, err := m.coll.InsertOne(ctx, c); err != nil {
		if mgo.IsDuplicateKeyError(err) {
			return commodity.Commodity{}, ErrAlreadyExists
		}
		return commodity.Commodity{}, err
	}
	return c, nil
}

func (m *Mongo) Get(ctx context.Context, id commodity.ID) (commodity.Commodity, error) {
	var c commodity.Commodity
	err := m.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&c)
	if errors.Is(err, mgo.ErrNoDocuments) {
		return commodity.Commodity{}, ErrNotFound
	}
	return c, err
}

func (m *Mongo) List(ctx context.Context) ([]commodity.Commodity, error) {
	cur, err := m.coll.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "name", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []commodity.Commodity
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (m *Mongo) Update(ctx context.Context, id commodity.ID, u Update) (commodity.Commodity, error) {
	if u.IsEmpty() {
		return m.Get(ctx, id)
	}

	// Read-modify-write so we can validate the merged record before writing.
	cur, err := m.Get(ctx, id)
	if err != nil {
		return commodity.Commodity{}, err
	}
	next := u.Apply(cur)
	if err := next.Validate(); err != nil {
		return commodity.Commodity{}, err
	}
	next.UpdatedAt = m.now()

	set := bson.M{"updated_at": next.UpdatedAt}
	if u.Name != nil {
		set["name"] = next.Name
	}
	if u.UseValueDescription != nil {
		set["use_value.description"] = next.UseValue.Description
	}
	if u.UseValueUnit != nil {
		set["use_value.unit"] = next.UseValue.Unit
	}
	if u.ConcreteLabourKind != nil {
		set["concrete_labour.kind"] = next.ConcreteLabour.Kind
	}
	if u.ConcreteLabourDescription != nil {
		set["concrete_labour.description"] = next.ConcreteLabour.Description
	}
	if u.SNLTPerUnit != nil {
		set["snlt_per_unit"] = next.SNLTPerUnit
	}

	res, err := m.coll.UpdateByID(ctx, id, bson.M{"$set": set})
	if err != nil {
		if mgo.IsDuplicateKeyError(err) {
			return commodity.Commodity{}, ErrAlreadyExists
		}
		return commodity.Commodity{}, err
	}
	if res.MatchedCount == 0 {
		return commodity.Commodity{}, ErrNotFound
	}
	return next, nil
}

func (m *Mongo) Delete(ctx context.Context, id commodity.ID) error {
	res, err := m.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

