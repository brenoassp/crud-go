package repo

import (
	"context"

	"github.com/brenoassp/crud-go/domain"
)

// Clients represents the operations we use for
// creating and retrieving a client from a persistent storage
type Clients interface {
	CreateClient(ctx context.Context, client domain.Client) error
	GetClients(ctx context.Context, page int, pageSize int) (domain.GetClientsResponse, error)
	DeleteClient(ctx context.Context, id int) error
	UpdateClient(ctx context.Context, client domain.Client) error
}
