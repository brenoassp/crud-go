package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/brenoassp/crud-go/adapters/cryptography"
	"github.com/brenoassp/crud-go/adapters/log"
	"github.com/brenoassp/crud-go/adapters/messageBroker"
	"github.com/brenoassp/crud-go/adapters/repo"
	"github.com/brenoassp/crud-go/domain"
)

type Service struct {
	logger        log.Provider
	clientsRepo   repo.Clients
	messageBroker messageBroker.Provider
	crypto        cryptography.Provider
}

func NewService(
	logger log.Provider,
	clientsRepo repo.Clients,
	messageBroker messageBroker.Provider,
	cryptoClient cryptography.Provider,
) *Service {
	return &Service{
		logger:        logger,
		clientsRepo:   clientsRepo,
		messageBroker: messageBroker,
		crypto:        cryptoClient,
	}
}

func (s *Service) GetClients(ctx context.Context, page int, pageSize int) (domain.GetClientsResponse, error) {
	response, err := s.clientsRepo.GetClients(ctx, page, pageSize)
	if err != nil {
		return response, fmt.Errorf("failed to get clients: %w", err)
	}
	return response, nil
}

func (s *Service) CreateClientEvent(ctx context.Context, event domain.CreateClientEvent) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal createClientEvent: %w", err)
	}

	encryptedData, err := s.crypto.Encrypt(eventBytes)
	if err != nil {
		return fmt.Errorf("failed to encrypt user data: %w", err)
	}

	err = s.messageBroker.Publish(
		ctx,
		"",        // exchange
		"clients", // routing key
		false,     // mandatory
		false,     // immediate
		messageBroker.Message{
			ContentType: "application/json",
			Body:        encryptedData,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

func (s *Service) CreateClientFromEvent(ctx context.Context, event domain.CreateClientEvent) error {
	format := "01-02-2006"
	birthdateTime, err := time.Parse(format, event.Birthdate)
	if err != nil {
		return fmt.Errorf("failed to parse birthdate: %w", err)
	}

	client := domain.Client{
		Name:      &event.Name,
		Surname:   &event.Surname,
		Contact:   &event.Contact,
		Address:   &event.Address,
		Birthdate: &birthdateTime,
		Cpf:       &event.Cpf,
	}

	err = s.clientsRepo.CreateClient(ctx, client)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) DeleteClient(ctx context.Context, id int) error {
	err := s.clientsRepo.DeleteClient(ctx, id)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateClient(ctx context.Context, client domain.Client) error {
	err := s.clientsRepo.UpdateClient(ctx, client)
	if err != nil {
		return err
	}
	return nil
}
