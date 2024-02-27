package pgrepo

import (
	"context"

	"github.com/brenoassp/crud-go/domain"
	"github.com/vingarcia/ksql"
)

func createClient(ctx context.Context, db ksql.Provider, client domain.Client) error {
	err := db.Insert(ctx, domain.ClientsTable, &client)
	if err != nil {
		return domain.InternalErr("unexpected error when saving client", map[string]interface{}{
			"client": client,
			"error":  err.Error(),
		})
	}
	return nil
}

func getClients(ctx context.Context, db ksql.Provider, page int, pageSize int) (domain.GetClientsResponse, error) {
	var response domain.GetClientsResponse = domain.GetClientsResponse{
		Page:     page,
		PageSize: pageSize,
		Total:    0,
		Data:     []domain.Client{},
	}

	var count []struct {
		Count int `ksql:"count"`
	}
	err := db.Query(ctx, &count, `SELECT COUNT(*) FROM clients`)
	if err != nil {
		return response, domain.InternalErr("unexpected error when fetching clients count", map[string]interface{}{
			"error": err.Error(),
		})
	}
	response.Total = count[0].Count

	err = db.Query(ctx, &response.Data, `SELECT * FROM clients LIMIT $1 OFFSET $2`, pageSize, (page-1)*pageSize)
	if err != nil {
		return response, domain.InternalErr("unexpected error when fetching clients", map[string]interface{}{
			"page":     page,
			"pageSize": pageSize,
			"error":    err.Error(),
		})
	}

	return response, nil
}

func deleteClient(ctx context.Context, db ksql.Provider, id int) error {
	client := domain.Client{}
	err := db.QueryOne(ctx, &client, `SELECT * FROM clients WHERE id = $1`, id)
	if err != nil && err != ksql.ErrRecordNotFound {
		return domain.InternalErr("unexpected error when fetching client to delete", map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		})
	}

	if err == ksql.ErrRecordNotFound {
		return domain.NotFoundErr("client not found", map[string]interface{}{
			"id": id,
		})
	}

	_, err = db.Exec(ctx, `DELETE FROM clients WHERE id = $1`, id)
	if err != nil {
		return domain.InternalErr("unexpected error when deleting client", map[string]interface{}{
			"id":    id,
			"error": err.Error(),
		})
	}
	return nil
}

func updateClient(ctx context.Context, db ksql.Provider, client domain.Client) error {
	oldClient := domain.Client{}
	err := db.QueryOne(ctx, &oldClient, `SELECT * FROM clients WHERE id = $1`, client.ID)
	if err != nil && err != ksql.ErrRecordNotFound {
		return domain.InternalErr("unexpected error when fetching client to update", map[string]interface{}{
			"id":    client.ID,
			"error": err.Error(),
		})
	}

	if err == ksql.ErrRecordNotFound {
		return domain.NotFoundErr("client not found", map[string]interface{}{
			"id": client.ID,
		})
	}

	err = db.Patch(ctx, domain.ClientsTable, client)
	if err != nil {
		return domain.InternalErr("unexpected error when updating client", map[string]interface{}{
			"client": client,
			"error":  err.Error(),
		})
	}
	return nil
}
