package main

import (
	"CRUD_Hexagonal/api/product"
	_ "CRUD_Hexagonal/api/product"
	"CRUD_Hexagonal/infrastructure"
	storeRepo "CRUD_Hexagonal/repository/product"
	storeServ "CRUD_Hexagonal/service/product"
	"context"
	"errors"
	otelfiber "github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	middlewareLogger "github.com/gofiber/fiber/v2/middleware/logger"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/viper"

	"golang.org/x/exp/slog"
	"log"
	"os"
)

func main() {
	// Initialize Viper
	viper.AutomaticEnv() // Read environment variables
	ctx := context.Background()

	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelInfo,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			return a
		},
	}).WithAttrs([]slog.Attr{
		slog.String("service", os.Getenv("OTEL_SERVICE_NAME")),
		slog.String("with-release", "v1.0.0"),
	})
	logger := slog.New(logHandler)
	slog.SetDefault(logger)

	// Set up OpenTelemetry.
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	otelCollector := os.Getenv("OTEL_COLLECTOR")
	serviceVersion := "0.1.0"
	otelShutdown, err := infrastructure.SetupOTelSDK(ctx, otelCollector, serviceName, serviceVersion, os.Getenv("OTEL_ENV"))
	if err != nil {
		log.Fatalf("failed to initialize OTel SDK: %v", err)
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()

	// init mongo
	mongo := infrastructure.NewMongo(ctx, os.Getenv("MONGO_DSN"), os.Getenv("MONGO_DB_NAME"))
	mongo = mongo.Connect()

	// store
	storeRepository := storeRepo.NewstoreRepository(mongo.Client, mongo.DB, "products")
	storeService := storeServ.NewStoreService(storeRepository)
	handler := product.NewStoreHandler(storeService)

	app := fiber.New()

	app.Use(middlewareLogger.New(middlewareLogger.Config{
		Format: "[${time}] ${ip}  ${status} - ${latency} ${method} ${path}\n",
	}))

	app.Use(otelfiber.Middleware())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Get("/product/:id", handler.Get)
	app.Get("/product", handler.GetAll)
	app.Post("/product", handler.Create)
	app.Put("/product/:id", handler.Update)
	app.Delete("/product/:id", handler.Delete)

	app.Listen(":3000")
}
