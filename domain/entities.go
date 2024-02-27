package domain

import "time"

type Client struct {
	ID        int        `ksql:"id"`
	Name      *string    `ksql:"name"`
	Surname   *string    `ksql:"surname"`
	Contact   *string    `ksql:"contact"`
	Address   *string    `ksql:"address"`
	Birthdate *time.Time `ksql:"birthdate"`
	Cpf       *string    `ksql:"cpf"`

	CreatedAt *time.Time `ksql:"created_at,timeNowUTC/skipUpdates"`
	UpdatedAt *time.Time `ksql:"updated_at,timeNowUTC"`
}

type CreateClientEvent struct {
	Name      string `json:"name"`
	Surname   string `json:"surname"`
	Contact   string `json:"contact"`
	Address   string `json:"address"`
	Birthdate string `json:"birthdate"`
	Cpf       string `json:"cpf"`
}

type GetClientsResponse struct {
	Data     []Client `json:"data"`
	Total    int      `json:"total"`
	Page     int      `json:"page"`
	PageSize int      `json:"page_size"`
}
