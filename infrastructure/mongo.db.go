package infrastructure

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
	"log/slog"
)

type AdapterMongo struct {
	MongoDSN string
	ctx      context.Context
	Client   *mongo.Client
	DB       string
}

func NewMongo(ctx context.Context, mongoDSN string, db string) *AdapterMongo {
	return &AdapterMongo{ctx: ctx, MongoDSN: mongoDSN, DB: db}
}

func (m *AdapterMongo) Connect() *AdapterMongo {
	clientOptions := options.Client()
	clientOptions.Monitor = otelmongo.NewMonitor()
	clientOptions.ApplyURI(m.MongoDSN)

	client, err := mongo.Connect(m.ctx, clientOptions)
	if err != nil {
		slog.ErrorContext(m.ctx, "error connect to mongo", slog.Any("err ", err))
	}

	slog.InfoContext(m.ctx, "Mongodb connected.")

	return &AdapterMongo{
		Client: client,
		DB:     m.DB,
	}

}
