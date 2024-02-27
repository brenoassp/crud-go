package main

import (
	"context"
	"os"
	"time"

	"github.com/brenoassp/crud-go/adapters/cryptography/crypto"
	"github.com/brenoassp/crud-go/adapters/log"
	"github.com/brenoassp/crud-go/adapters/log/jsonlogs"
	"github.com/brenoassp/crud-go/adapters/messageBroker"
	"github.com/brenoassp/crud-go/adapters/messageBroker/rabbitmq"
	"github.com/brenoassp/crud-go/adapters/repo"
	pgrepo "github.com/brenoassp/crud-go/adapters/repo/pg_repo"
	"github.com/brenoassp/crud-go/cmd/api/clientsctrl"
	"github.com/brenoassp/crud-go/cmd/api/middlewares"
	"github.com/brenoassp/crud-go/domain"
	"github.com/brenoassp/crud-go/domain/clients"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	ctx := context.Background()
	godotenv.Load("config.env")

	port := os.Getenv("PORT")
	logLevel := os.Getenv("LOG_LEVEL")
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	dbURL := os.Getenv("DATABASE_URL")
	encryptionKey := os.Getenv("CRYPTO_KEY")

	logger := jsonlogs.New(logLevel, domain.GetCtxValues)

	err := tryMigrate(ctx, dbURL)
	if err == errMigrationRunning {
		logger.Error(ctx, "error-migration-already-running", log.Body{
			"error": "a migration is being run by a different instance, restarting...",
		})
		time.Sleep(500 * time.Millisecond)
		return
	}
	if err != nil {
		logger.Fatal(ctx, "error-running-migrations", log.Body{
			"error": err.Error(),
		})
	}

	var clientsRepo repo.Clients
	clientsRepo, err = pgrepo.New(ctx, dbURL)
	if err != nil {
		logger.Fatal(ctx, "unable to start database", log.Body{
			"db_url": dbURL,
			"error":  err.Error(),
		})
	}

	var rabbitmqClient messageBroker.Provider = rabbitmq.New(rabbitmqURL)

	cryptoClient := crypto.NewClient(encryptionKey)

	clientsService := clients.NewService(logger, clientsRepo, rabbitmqClient, cryptoClient)

	clientsController := clientsctrl.NewController(clientsService)

	app := fiber.New()

	app.Use(middlewares.HandleError(logger))

	app.Post("/clients", clientsController.CreateClient)
	app.Get("/clients", clientsController.GetClients)
	app.Patch("/clients/:id", clientsController.UpdateClient)
	app.Delete("/clients/:id", clientsController.DeleteClient)

	logger.Info(ctx, "server-starting-up", log.Body{
		"port": port,
	})
	if err := app.Listen(":" + port); err != nil {
		logger.Error(ctx, "server-stopped-with-an-error", log.Body{
			"error": err.Error(),
		})
	}
}
