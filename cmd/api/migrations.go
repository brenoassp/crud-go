package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/brenoassp/crud-go/migrations"
	_ "github.com/jackc/pgx/v4/stdlib"
	migrate "github.com/rubenv/sql-migrate"
)

var errMigrationRunning error = fmt.Errorf("a migration is already being executed by a different process")

func tryMigrate(ctx context.Context, dbURL string) (err error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("tryMigrate: unable connect to db: %w", err)
	}
	defer db.Close()

	// When using postgres advisory locks it is very important the lock and unlock
	// operations happen on the same connection:
	conn, err := db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("tryMigrate: unable to obtain a db connection: %w", err)
	}

	var lockSuccess bool
	err = conn.QueryRowContext(ctx, "SELECT pg_try_advisory_lock(4242)").Scan(&lockSuccess)
	if err != nil {
		return fmt.Errorf("tryMigrate: unexpected error when obtaining a postgres lock: %w", err)
	}

	if !lockSuccess {
		return errMigrationRunning
	}

	defer func() {
		_, execErr := conn.ExecContext(ctx, "SELECT pg_advisory_unlock(4242)")
		if execErr != nil {
			err = errors.Join(err, fmt.Errorf("tryMigrate: unexpected error releasing db lock: %w", execErr))
		}
	}()

	migrator := migrate.MigrationSet{
		TableName: "migrations",
	}

	n, err := migrator.Exec(db, "postgres", &migrate.HttpFileSystemMigrationSource{
		FileSystem: http.FS(migrations.Dir),
	}, migrate.Up)
	if err != nil {
		return fmt.Errorf("a migration failed after %d successful migrations: %w", n, err)
	}

	if n > 0 {
		fmt.Printf("successfully executed %d migrations\n", n)
	}

	return nil
}
