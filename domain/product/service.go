package product

import (
	"CRUD_Hexagonal/utils"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProductInterface , interface defining all CRUD operations
type ProductInterface interface {
	Find(ctx context.Context, id string) (*Product, error)
	Store(ctx context.Context, product *Product) (*Product, error)
	Update(ctx context.Context, dataStore *Product) error
	FindAll(ctx context.Context, filter Filter) ([]*Product, *utils.Pagination, error)
	Delete(ctx context.Context, code string) error
	DeleteById(ctx context.Context, id string) error
}

// Repository , interface acting like a port for the database implementation
type Repository interface {
	Find(ctx context.Context, id string) (*Product, error)
	Store(ctx context.Context, dataStore *Product) (primitive.ObjectID, error)
	Update(ctx context.Context, dataStore *Product) error
	FindAll(ctx context.Context, filter Filter) ([]*Product, *utils.Pagination, error)
	Delete(ctx context.Context, code string) error
	DeleteById(ctx context.Context, id string) error
}
