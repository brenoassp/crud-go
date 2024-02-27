package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/brenoassp/crud-go/adapters/cryptography"
	"github.com/brenoassp/crud-go/adapters/cryptography/crypto"
	"github.com/brenoassp/crud-go/adapters/log"
	"github.com/brenoassp/crud-go/adapters/log/jsonlogs"
	"github.com/brenoassp/crud-go/adapters/messageBroker"
	"github.com/brenoassp/crud-go/adapters/messageBroker/rabbitmq"
	"github.com/brenoassp/crud-go/adapters/repo"
	pgrepo "github.com/brenoassp/crud-go/adapters/repo/pg_repo"
	"github.com/brenoassp/crud-go/domain"
	"github.com/brenoassp/crud-go/domain/clients"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()

	// Load the configuration file
	godotenv.Load("config.env")
	dbURL := os.Getenv("DATABASE_URL")
	logger := jsonlogs.New("INFO")
	rabbitmqURL := os.Getenv("RABBITMQ_URL")
	encryptionKey := os.Getenv("CRYPTO_KEY")

	var clientsRepo repo.Clients
	clientsRepo, err := pgrepo.New(ctx, dbURL)
	if err != nil {
		logger.Fatal(ctx, "unable to start database", log.Body{
			"db_url": dbURL,
			"error":  err.Error(),
		})
	}

	var rabbitmqClient messageBroker.Provider = rabbitmq.New(rabbitmqURL)
	cryptoClient := crypto.NewClient(encryptionKey)

	clientsService := clients.NewService(logger, clientsRepo, rabbitmqClient, cryptoClient)

	go Consume(ctx, rabbitmqClient, cryptoClient, clientsService)

	<-ctx.Done()
}

func Consume(
	ctx context.Context,
	rabbitmqClient messageBroker.Provider,
	cryptoClient cryptography.Provider,
	clientsService *clients.Service,
) {
	consumer, err := rabbitmqClient.Consume(ctx, "", "clients", "direct")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		select {
		case <-ctx.Done():
			fmt.Println("consumer done")
			consumer.Close()
			return
		case msg := <-consumer.Message:
			unencryptedEvent, err := cryptoClient.Decrypt(msg.Body)
			if err != nil {
				fmt.Println(err)
				msg.Nack(false, true)
				continue
			}

			var event domain.CreateClientEvent
			err = json.Unmarshal(unencryptedEvent, &event)
			if err != nil {
				fmt.Println(err)
				msg.Nack(false, true)
				continue
			}

			err = clientsService.CreateClientFromEvent(ctx, event)
			if err != nil {
				fmt.Printf("failed to create client: %v", err)
				msg.Nack(false, true)
				continue
			}

			err = msg.Ack(true)
			if err != nil {
				fmt.Printf("failed to ack message: %v", err)
				continue
			}
		}
	}
}
