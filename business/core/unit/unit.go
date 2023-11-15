package unit

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/nhaancs/bhms/business/core/block"
	"github.com/nhaancs/bhms/business/core/floor"
	"github.com/nhaancs/bhms/business/core/property"
	"github.com/nhaancs/bhms/business/data/transaction"
	"github.com/nhaancs/bhms/foundation/logger"
	"time"
)

var (
	ErrNotFound      = errors.New("unit not found")
	ErrLimitExceeded = errors.New("max number of units exceeded")
	ErrBlockNotFound = errors.New("block not found")
	ErrFloorNotFound = errors.New("floor not found")
)

type Storer interface {
	Create(ctx context.Context, core Unit) error
	BatchCreate(ctx context.Context, cores []Unit) error
	DeleteByPropertyID(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, core Unit) error
	QueryByID(ctx context.Context, id uuid.UUID) (Unit, error)
	QueryByFloorID(ctx context.Context, id uuid.UUID) ([]Unit, error)
	QueryByPropertyID(ctx context.Context, id uuid.UUID) ([]Unit, error)
	ExecuteUnderTransaction(tx transaction.Transaction) (Storer, error)
}

type Core struct {
	store    Storer
	log      *logger.Logger
	property *property.Core
	block    *block.Core
	floor    *floor.Core
}

func NewCore(
	log *logger.Logger,
	store Storer,
	property *property.Core,
	block *block.Core,
	floor *floor.Core,
) *Core {
	return &Core{
		store:    store,
		log:      log,
		property: property,
		block:    block,
		floor:    floor,
	}
}

func (c *Core) Create(ctx context.Context, core NewUnit) (Unit, error) {
	unts, err := c.store.QueryByPropertyID(ctx, core.PropertyID)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return Unit{}, fmt.Errorf("QueryByPropertyID: id[%s]: %w", core.PropertyID, err)
	}
	if len(unts) >= 1000 {
		return Unit{}, ErrLimitExceeded
	}

	_, err = c.block.QueryByID(ctx, core.BlockID)
	if err != nil {
		switch {
		case errors.Is(err, block.ErrNotFound):
			return Unit{}, ErrBlockNotFound
		default:
			return Unit{}, fmt.Errorf("querybyid: block id[%s]: %w", core.BlockID, err)
		}
	}

	_, err = c.floor.QueryByID(ctx, core.FloorID)
	if err != nil {
		switch {
		case errors.Is(err, floor.ErrNotFound):
			return Unit{}, ErrFloorNotFound
		default:
			return Unit{}, fmt.Errorf("querybyid: floor id[%s]: %w", core.FloorID, err)
		}
	}

	now := time.Now()
	unt := Unit{
		ID:         core.ID,
		Name:       core.Name,
		PropertyID: core.PropertyID,
		BlockID:    core.BlockID,
		FloorID:    core.FloorID,
		Status:     StatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := c.store.Create(ctx, unt); err != nil {
		return Unit{}, fmt.Errorf("create: %w", err)
	}

	return unt, nil
}

func (c *Core) BatchCreate(ctx context.Context, cores []NewUnit) ([]Unit, error) {
	now := time.Now()
	unts := make([]Unit, len(cores))

	for i := range cores {
		unts[i] = Unit{
			ID:         uuid.New(),
			Name:       cores[i].Name,
			PropertyID: cores[i].PropertyID,
			BlockID:    cores[i].BlockID,
			FloorID:    cores[i].FloorID,
			Status:     StatusActive,
			CreatedAt:  now,
			UpdatedAt:  now,
		}
	}

	if err := c.store.BatchCreate(ctx, unts); err != nil {
		return nil, fmt.Errorf("batch create: %w", err)
	}

	return unts, nil
}

func (c *Core) Update(ctx context.Context, o Unit, n UpdateUnit) (Unit, error) {
	if n.Name != nil {
		o.Name = *n.Name
	}

	o.UpdatedAt = time.Now()

	if err := c.store.Update(ctx, o); err != nil {
		return Unit{}, fmt.Errorf("update: %w", err)
	}

	return o, nil
}

func (c *Core) Delete(ctx context.Context, core Unit) (Unit, error) {
	core.UpdatedAt = time.Now()
	core.Status = StatusDeleted

	if err := c.store.Update(ctx, core); err != nil {
		return Unit{}, fmt.Errorf("update: %w", err)
	}

	return core, nil
}

func (c *Core) DeleteByPropertyID(ctx context.Context, propertyID uuid.UUID) error {
	if err := c.store.DeleteByPropertyID(ctx, propertyID); err != nil {
		return fmt.Errorf("DeleteByPropertyID: %w", err)
	}
	return nil
}

func (c *Core) QueryByID(ctx context.Context, id uuid.UUID) (Unit, error) {
	unt, err := c.store.QueryByID(ctx, id)
	if err != nil {
		return Unit{}, fmt.Errorf("query: id[%s]: %w", id, err)
	}

	return unt, nil
}

func (c *Core) QueryByFloorID(ctx context.Context, id uuid.UUID) ([]Unit, error) {
	unts, err := c.store.QueryByFloorID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("query: block id[%s]: %w", id.String(), err)
	}

	return unts, nil
}

func (c *Core) QueryByPropertyID(ctx context.Context, id uuid.UUID) ([]Unit, error) {
	unts, err := c.store.QueryByPropertyID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("query: property id[%s]: %w", id.String(), err)
	}

	return unts, nil
}

// ExecuteUnderTransaction constructs a new Core value that will use the
// specified transaction in any store related calls.
func (c *Core) ExecuteUnderTransaction(tx transaction.Transaction) (*Core, error) {
	store, err := c.store.ExecuteUnderTransaction(tx)
	if err != nil {
		return nil, err
	}

	c = &Core{
		store: store,
		log:   c.log,
	}

	return c, nil
}
