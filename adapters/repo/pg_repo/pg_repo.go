package pgrepo

import (
	"context"

	"github.com/brenoassp/crud-go/adapters/log"
	"github.com/brenoassp/crud-go/domain"
	"github.com/vingarcia/ksql"
	"github.com/vingarcia/ksql/adapters/kpgx"
)

// Repo implements the repo.Clients interface by using the ksql database.
type Repo struct {
	db ksql.Provider
}

// New instantiates a new Repo
func New(ctx context.Context, postgresURL string) (Repo, error) {
	db, err := kpgx.New(ctx, postgresURL, ksql.Config{})
	if err != nil {
		return Repo{}, domain.InternalErr("unable to start database", log.Body{
			"error": err.Error(),
		})
	}

	return Repo{
		db: db,
	}, nil
}

// CreateClient implements the repo.Clients interface
func (u Repo) CreateClient(ctx context.Context, client domain.Client) error {
	return createClient(ctx, u.db, client)
}

// CreateClient implements the repo.Clients interface
func (u Repo) GetClients(ctx context.Context, page int, pageSize int) (domain.GetClientsResponse, error) {
	return getClients(ctx, u.db, page, pageSize)
}

func (u Repo) DeleteClient(ctx context.Context, id int) error {
	return deleteClient(ctx, u.db, id)
}

func (u Repo) UpdateClient(ctx context.Context, client domain.Client) error {
	return updateClient(ctx, u.db, client)
}
