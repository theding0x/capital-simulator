// Package store defines the persistence boundary for commodity-service.
// The Store interface is the only seam between the HTTP/domain layer and
// the underlying database (MySQL in production, in-memory in tests).
package store

import (
	"context"
	"errors"

	"github.com/theding0x/capital-simulator/services/commodity-service/internal/commodity"
)

// Sentinel errors that callers can branch on with errors.Is.
var (
	ErrNotFound      = errors.New("commodity: not found")
	ErrAlreadyExists = errors.New("commodity: already exists")
)

// Update is a partial-update payload. Non-nil fields are applied; nil fields
// are left untouched. Modeling it this way keeps PATCH semantics explicit
// and avoids the "did the caller mean empty or did they mean unchanged?"
// ambiguity of plain structs.
type Update struct {
	Name                      *string
	UseValueDescription       *string
	UseValueUnit              *string
	ConcreteLabourKind        *string
	ConcreteLabourDescription *string
	SNLTPerUnit               *commodity.LabourMinutes
}

// IsEmpty reports whether u carries no field updates.
func (u Update) IsEmpty() bool {
	return u.Name == nil &&
		u.UseValueDescription == nil &&
		u.UseValueUnit == nil &&
		u.ConcreteLabourKind == nil &&
		u.ConcreteLabourDescription == nil &&
		u.SNLTPerUnit == nil
}

// Apply returns a copy of c with u's non-nil fields written through.
func (u Update) Apply(c commodity.Commodity) commodity.Commodity {
	out := c
	if u.Name != nil {
		out.Name = *u.Name
	}
	if u.UseValueDescription != nil {
		out.UseValue.Description = *u.UseValueDescription
	}
	if u.UseValueUnit != nil {
		out.UseValue.Unit = *u.UseValueUnit
	}
	if u.ConcreteLabourKind != nil {
		out.ConcreteLabour.Kind = *u.ConcreteLabourKind
	}
	if u.ConcreteLabourDescription != nil {
		out.ConcreteLabour.Description = *u.ConcreteLabourDescription
	}
	if u.SNLTPerUnit != nil {
		out.SNLTPerUnit = *u.SNLTPerUnit
	}
	return out
}

// Store is the persistence contract for Commodity records.
type Store interface {
	Create(ctx context.Context, c commodity.Commodity) (commodity.Commodity, error)
	Get(ctx context.Context, id commodity.ID) (commodity.Commodity, error)
	List(ctx context.Context) ([]commodity.Commodity, error)
	Update(ctx context.Context, id commodity.ID, u Update) (commodity.Commodity, error)
	Delete(ctx context.Context, id commodity.ID) error
}
