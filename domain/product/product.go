package product

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ID        primitive.ObjectID `json:"product_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"product_name" bson:"product_name"`
	Stock     int64              `json:"stock" bson:"stock"`
	CreatedAt int64              `json:"created_at"`
	UpdatedAt int64              `json:"updated_at"`
	DeletedAt int64              `json:"deleted_at"`
}

type Filter struct {
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
	Keyword   string `json:"keyword"`
}
