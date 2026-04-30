package store

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	mgo "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/theding0x/capital-simulator/services/market-service/internal/market"
)

const (
	colOwners    = "owners"
	colOffers    = "offers"
	colExchanges = "exchanges"
	colConfig    = "market_config"
	colPrices    = "prices"

	// Singleton document keys stored in market_config.
	keyUniversalEquivalent = "universal_equivalent"
	keyMoneyCommodity      = "money_commodity"
)

// Mongo is a MongoDB-backed Store.
type Mongo struct {
	owners    *mgo.Collection
	offers    *mgo.Collection
	exchanges *mgo.Collection
	config    *mgo.Collection
	prices    *mgo.Collection
	now       func() time.Time
}

// NewMongo returns a Store backed by the given database and ensures necessary
// indexes exist.
func NewMongo(ctx context.Context, db *mgo.Database) (*Mongo, error) {
	s := &Mongo{
		owners:    db.Collection(colOwners),
		offers:    db.Collection(colOffers),
		exchanges: db.Collection(colExchanges),
		config:    db.Collection(colConfig),
		prices:    db.Collection(colPrices),
		now:       time.Now,
	}
	// Unique index on prices keyed by commodity_id.
	if _, err := s.prices.Indexes().CreateOne(ctx, mgo.IndexModel{
		Keys:    bson.D{{Key: "commodity_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, err
	}
	return s, nil
}

// --- Owner ------------------------------------------------------------------

func (m *Mongo) CreateOwner(ctx context.Context, o market.Owner) (market.Owner, error) {
	if err := o.Validate(); err != nil {
		return market.Owner{}, err
	}
	if o.ID.IsZero() {
		o.ID = market.NewOwnerID()
	}
	now := m.now()
	o.CreatedAt = now
	o.UpdatedAt = now
	if _, err := m.owners.InsertOne(ctx, o); err != nil {
		if mgo.IsDuplicateKeyError(err) {
			return market.Owner{}, ErrAlreadyExists
		}
		return market.Owner{}, err
	}
	return o, nil
}

func (m *Mongo) GetOwner(ctx context.Context, id market.OwnerID) (market.Owner, error) {
	var o market.Owner
	if err := m.owners.FindOne(ctx, bson.M{"_id": id}).Decode(&o); err != nil {
		if errors.Is(err, mgo.ErrNoDocuments) {
			return market.Owner{}, ErrNotFound
		}
		return market.Owner{}, err
	}
	return o, nil
}

func (m *Mongo) ListOwners(ctx context.Context) ([]market.Owner, error) {
	cur, err := m.owners.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []market.Owner
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// --- Offer ------------------------------------------------------------------

func (m *Mongo) CreateOffer(ctx context.Context, o market.Offer) (market.Offer, error) {
	if err := o.Validate(); err != nil {
		return market.Offer{}, err
	}
	if o.ID.IsZero() {
		o.ID = market.NewOfferID()
	}
	o.CreatedAt = m.now()
	if _, err := m.offers.InsertOne(ctx, o); err != nil {
		if mgo.IsDuplicateKeyError(err) {
			return market.Offer{}, ErrAlreadyExists
		}
		return market.Offer{}, err
	}
	return o, nil
}

func (m *Mongo) GetOffer(ctx context.Context, id market.OfferID) (market.Offer, error) {
	var o market.Offer
	if err := m.offers.FindOne(ctx, bson.M{"_id": id}).Decode(&o); err != nil {
		if errors.Is(err, mgo.ErrNoDocuments) {
			return market.Offer{}, ErrNotFound
		}
		return market.Offer{}, err
	}
	return o, nil
}

func (m *Mongo) ListOffers(ctx context.Context) ([]market.Offer, error) {
	cur, err := m.offers.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []market.Offer
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

func (m *Mongo) DeleteOffer(ctx context.Context, id market.OfferID) error {
	res, err := m.offers.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// --- Exchange ---------------------------------------------------------------

func (m *Mongo) CreateExchange(ctx context.Context, e market.Exchange) (market.Exchange, error) {
	if err := e.Validate(); err != nil {
		return market.Exchange{}, err
	}
	if e.ID.IsZero() {
		e.ID = market.NewExchangeID()
	}
	e.CreatedAt = m.now()
	if _, err := m.exchanges.InsertOne(ctx, e); err != nil {
		if mgo.IsDuplicateKeyError(err) {
			return market.Exchange{}, ErrAlreadyExists
		}
		return market.Exchange{}, err
	}
	return e, nil
}

func (m *Mongo) GetExchange(ctx context.Context, id market.ExchangeID) (market.Exchange, error) {
	var e market.Exchange
	if err := m.exchanges.FindOne(ctx, bson.M{"_id": id}).Decode(&e); err != nil {
		if errors.Is(err, mgo.ErrNoDocuments) {
			return market.Exchange{}, ErrNotFound
		}
		return market.Exchange{}, err
	}
	return e, nil
}

func (m *Mongo) ListExchanges(ctx context.Context) ([]market.Exchange, error) {
	cur, err := m.exchanges.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []market.Exchange
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// --- UniversalEquivalent (singleton stored in market_config) ----------------

func (m *Mongo) SetUniversalEquivalent(ctx context.Context, ue market.UniversalEquivalent) (market.UniversalEquivalent, error) {
	// Idempotent: read first and check commodity_id.
	existing, err := m.GetUniversalEquivalent(ctx)
	if err == nil && existing.CommodityID == ue.CommodityID {
		return existing, nil
	}
	if ue.SetAt.IsZero() {
		ue.SetAt = m.now()
	}
	doc := bson.M{"_id": keyUniversalEquivalent, "commodity_id": ue.CommodityID, "set_at": ue.SetAt}
	_, err = m.config.ReplaceOne(ctx, bson.M{"_id": keyUniversalEquivalent}, doc, options.Replace().SetUpsert(true))
	if err != nil {
		return market.UniversalEquivalent{}, err
	}
	return ue, nil
}

func (m *Mongo) GetUniversalEquivalent(ctx context.Context) (market.UniversalEquivalent, error) {
	var doc struct {
		CommodityID market.CommodityID `bson:"commodity_id"`
		SetAt       time.Time          `bson:"set_at"`
	}
	if err := m.config.FindOne(ctx, bson.M{"_id": keyUniversalEquivalent}).Decode(&doc); err != nil {
		if errors.Is(err, mgo.ErrNoDocuments) {
			return market.UniversalEquivalent{}, ErrNotFound
		}
		return market.UniversalEquivalent{}, err
	}
	return market.UniversalEquivalent{CommodityID: doc.CommodityID, SetAt: doc.SetAt}, nil
}

// --- MoneyCommodity (singleton stored in market_config) ---------------------

func (m *Mongo) SetMoneyCommodity(ctx context.Context, mc market.MoneyCommodity) (market.MoneyCommodity, error) {
	if mc.CreatedAt.IsZero() {
		mc.CreatedAt = m.now()
	}
	doc := bson.M{"_id": keyMoneyCommodity, "commodity_id": mc.CommodityID, "created_at": mc.CreatedAt}
	_, err := m.config.ReplaceOne(ctx, bson.M{"_id": keyMoneyCommodity}, doc, options.Replace().SetUpsert(true))
	if err != nil {
		return market.MoneyCommodity{}, err
	}
	return mc, nil
}

func (m *Mongo) GetMoneyCommodity(ctx context.Context) (market.MoneyCommodity, error) {
	var doc struct {
		CommodityID market.CommodityID `bson:"commodity_id"`
		CreatedAt   time.Time          `bson:"created_at"`
	}
	if err := m.config.FindOne(ctx, bson.M{"_id": keyMoneyCommodity}).Decode(&doc); err != nil {
		if errors.Is(err, mgo.ErrNoDocuments) {
			return market.MoneyCommodity{}, ErrNotFound
		}
		return market.MoneyCommodity{}, err
	}
	return market.MoneyCommodity{CommodityID: doc.CommodityID, CreatedAt: doc.CreatedAt}, nil
}

// --- Price ------------------------------------------------------------------

func (m *Mongo) SetPrice(ctx context.Context, p market.Price) (market.Price, error) {
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = m.now()
	}
	doc := bson.M{
		"commodity_id":       p.CommodityID,
		"money_commodity_id": p.MoneyCommodityID,
		"amount":             p.Amount,
		"updated_at":         p.UpdatedAt,
	}
	_, err := m.prices.ReplaceOne(
		ctx,
		bson.M{"commodity_id": p.CommodityID},
		doc,
		options.Replace().SetUpsert(true),
	)
	if err != nil {
		return market.Price{}, err
	}
	return p, nil
}

func (m *Mongo) GetPrice(ctx context.Context, commodityID market.CommodityID) (market.Price, error) {
	var p market.Price
	if err := m.prices.FindOne(ctx, bson.M{"commodity_id": commodityID}).Decode(&p); err != nil {
		if errors.Is(err, mgo.ErrNoDocuments) {
			return market.Price{}, ErrNotFound
		}
		return market.Price{}, err
	}
	return p, nil
}

func (m *Mongo) ListPrices(ctx context.Context) ([]market.Price, error) {
	cur, err := m.prices.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "commodity_id", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []market.Price
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}
