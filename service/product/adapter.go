package product

import (
	"CRUD_Hexagonal/domain/product"
	_ "CRUD_Hexagonal/domain/product"
	"CRUD_Hexagonal/infrastructure"
	"CRUD_Hexagonal/utils"
	"context"
	"time"
)

type adapter struct {
	storeRepo product.Repository
}

func NewStoreService(storeRepo product.Repository) product.ProductInterface {
	return &adapter{storeRepo: storeRepo}
}

func (a adapter) Find(ctx context.Context, id string) (*product.Product, error) {
	// Tracing
	ctx, span := infrastructure.Tracer().Start(ctx, "service:store:Find")
	defer span.End()

	return a.storeRepo.Find(ctx, id)
}

func (a adapter) Store(ctx context.Context, product *product.Product) (*product.Product, error) {
	// Tracing
	ctx, span := infrastructure.Tracer().Start(ctx, "service:store:Store")
	defer span.End()

	product.CreatedAt = time.Now().UTC().Unix()

	insertID, err := a.storeRepo.Store(ctx, product)

	product.ID = insertID

	return product, err
}

func (a adapter) Update(ctx context.Context, store *product.Product) error {
	/// Tracing
	ctx, span := infrastructure.Tracer().Start(ctx, "service:store:Update")
	defer span.End()

	store.UpdatedAt = time.Now().UTC().Unix()

	return a.storeRepo.Update(ctx, store)
}

func (a adapter) FindAll(ctx context.Context, filter product.Filter) ([]*product.Product, *utils.Pagination, error) {
	// Tracing
	ctx, span := infrastructure.Tracer().Start(ctx, "service:store:FindAll")
	defer span.End()
	res, pagination, err := a.storeRepo.FindAll(ctx, filter)

	return res, pagination, err
}

func (a adapter) Delete(ctx context.Context, code string) error {
	//TODO implement me
	panic("implement me")
}

func (a adapter) DeleteById(ctx context.Context, id string) error {
	// Tracing
	ctx, span := infrastructure.Tracer().Start(ctx, "service:store:DeleteByID")
	defer span.End()

	err := a.storeRepo.DeleteById(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
