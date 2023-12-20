package product

import (
	"CRUD_Hexagonal/domain/product"
	_ "CRUD_Hexagonal/domain/product"
	"CRUD_Hexagonal/infrastructure"
	"CRUD_Hexagonal/utils"
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/exp/slog"
	"reflect"
)

type storeRepository struct {
	client     *mongo.Client
	db         string
	collection string
}

func NewstoreRepository(client *mongo.Client, db string, collection string) product.Repository {
	return &storeRepository{
		client:     client,
		db:         db,
		collection: collection,
	}
}

func (r *storeRepository) Find(ctx context.Context, id string) (*product.Product, error) {

	// Tracing
	ctx, span := infrastructure.Tracer().Start(ctx, "repository:store:Find")
	defer span.End()

	var storeData product.Product

	collection := r.client.Database(r.db).Collection(r.collection)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		slog.ErrorContext(ctx, "Collection Error", slog.Any("err ", err))
		return nil, err
	}
	filter := bson.D{{"_id", objectId}}
	err = collection.FindOne(ctx, filter).Decode(&storeData)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("error Finding a store")
		}
		slog.ErrorContext(ctx, "Collection Error", slog.Any("err ", err))
	}
	return &storeData, nil
}

func (r *storeRepository) FindAll(ctx context.Context, filter product.Filter) ([]*product.Product, *utils.Pagination, error) {
	// Tracing
	ctx, span := infrastructure.Tracer().Start(ctx, "repository:store:FindAll")
	defer span.End()

	collection := r.client.Database(r.db).Collection(r.collection)

	// Pagination
	var currentPage, limit int
	if filter.Limit <= 0 || filter.Page <= 0 {
		currentPage, limit = 1, 10
	} else {
		currentPage, limit = filter.Page, filter.Limit
	}
	skip := (currentPage - 1) * limit

	// Total Documents for Pagination
	totalDocuments, err := collection.CountDocuments(ctx, bson.D{})
	if err != nil {
		return nil, nil, err
	}

	pagination := utils.Pagination{
		Total:       int(totalDocuments),
		Limit:       limit,
		CurrentPage: currentPage,
	}

	// Apply Pagination
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))

	//filter latitude and longitude
	bsonFilter := bson.D{}

	if filter.Latitude != "" && filter.Longitude != "" {
		bsonFilter = append(bsonFilter, bson.E{
			Key:   "address.geo.latitude",
			Value: bson.D{{Key: "$eq", Value: filter.Latitude}},
		}, bson.E{
			Key:   "address.geo.longitude",
			Value: bson.D{{Key: "$eq", Value: filter.Longitude}},
		})
	}

	//filter by keyword
	if filter.Keyword != "" {
		bsonFilter = append(bsonFilter, bson.E{
			Key:   "name",
			Value: bson.D{{Key: "$regex", Value: primitive.Regex{Pattern: filter.Keyword, Options: "i"}}},
		})
	}

	var stores []*product.Product
	cur, err := collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil, errors.New("error finding store")
		}
		slog.ErrorContext(ctx, "Collection Error", slog.Any("err ", err))
		return nil, nil, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var elem product.Product
		if err := cur.Decode(&elem); err != nil {
			slog.ErrorContext(ctx, "Collection Error", slog.Any("err ", err))
			continue
		}
		stores = append(stores, &elem)
	}
	return stores, &pagination, nil
}

func (r *storeRepository) Store(ctx context.Context, dataStore *product.Product) (primitive.ObjectID, error) {
	// Tracing
	ctx, span := infrastructure.Tracer().Start(ctx, "repository:store:Store")
	defer span.End()

	collection := r.client.Database(r.db).Collection(r.collection)
	doInsert, err := collection.InsertOne(ctx, dataStore)
	if err != nil {
		slog.ErrorContext(ctx, "Error writing to repository", slog.Any("err ", err))
		return primitive.ObjectID{}, errors.New("error writing to repository")
	}

	return doInsert.InsertedID.(primitive.ObjectID), nil

}

func (r *storeRepository) Update(ctx context.Context, dataStore *product.Product) error {
	// Tracing
	ctx, span := infrastructure.Tracer().Start(ctx, "repository:store:Update")
	defer span.End()

	updatedStore := bson.D{}

	values := reflect.ValueOf(*dataStore)
	types := values.Type()
	for i := 0; i < values.NumField(); i++ {

		if types.Field(i).Name != "ID" && !utils.IsEmptyStruct(values.Field(i)) {
			updatedStore = append(updatedStore, primitive.E{Key: types.Field(i).Tag.Get("json"), Value: values.Field(i).Interface()})
		}
	}

	collection := r.client.Database(r.db).Collection(r.collection)
	_, err := collection.UpdateOne(
		ctx,
		bson.M{"_id": dataStore.ID},
		bson.D{
			{Key: "$set", Value: updatedStore},
		},
	)
	if err != nil {
		fmt.Println(err, "err")
		return err
	}

	return nil
}

func (r *storeRepository) Delete(ctx context.Context, code string) error {
	//TODO implement me
	panic("implement me")
}

func (r *storeRepository) DeleteById(ctx context.Context, id string) error {
	// Tracing
	ctx, span := infrastructure.Tracer().Start(ctx, "repository:store:DeleteById")
	defer span.End()

	if id == "" {
		return errors.New("ID is empty")
	}

	collection := r.client.Database(r.db).Collection(r.collection)

	objectID, err2 := primitive.ObjectIDFromHex(id)
	if err2 != nil {
		return err2
	}
	result := collection.FindOne(ctx, bson.M{"_id": objectID})

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return errors.New("store not found")
		}
		return result.Err()
	}

	_, err := collection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err != nil {
		return err
	}

	return nil

}
