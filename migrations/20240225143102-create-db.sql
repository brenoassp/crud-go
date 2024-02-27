-- +migrate Up
create table clients (
    id serial primary key,
    name varchar NOT NULL,
    surname varchar NOT NULL,
    contact varchar NOT NULL,
    address varchar NOT NULL,
    birthdate date NOT NULL,
    cpf varchar(11) NOT NULL,
    created_at timestamp NOT NULL DEFAULT NOW(),
    updated_at timestamp NOT NULL DEFAULT NOW()
);

-- +migrate Down
drop table clients;