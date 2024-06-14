package database

import (
	"birthday-service/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"sync"
)

type Postgres struct {
	Db     *pgxpool.Pool
	log    *slog.Logger
	Config *config.Config
}

var (
	pgInstance *Postgres
	pgOnce     sync.Once
)

func NewPG(ctx context.Context, connString string, log *slog.Logger, cfg *config.Config) (*Postgres, error) {
	var err error

	pgOnce.Do(func() {
		var db *pgxpool.Pool
		db, err = pgxpool.New(ctx, connString)

		if err != nil {
			err = fmt.Errorf("unable to create connection pool: %w", err)
			return
		}

		pgInstance = &Postgres{db, log, cfg}

		if err = CreateTables(ctx, db, log, cfg); err != nil {
			return
		}
	})

	if err != nil {
		return nil, err
	}
	return pgInstance, nil
}

func CreateTables(ctx context.Context, db *pgxpool.Pool, log *slog.Logger, cfg *config.Config) error {
	_, err := db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS Users (
	    id SERIAL PRIMARY KEY, 
	    email VARCHAR(100) UNIQUE NOT NULL, 
	    password VARCHAR(255) NOT NULL,
	    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP)
`)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	_, err = db.Exec(ctx, `CREATE TABLE IF NOT EXISTS Employees (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL, 
    birthday DATE NOT NULL)`)
	if err != nil {
		return fmt.Errorf("failed to create employee table: %w", err)
	}

	_, err = db.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS Subscriptions (
	id SERIAL PRIMARY KEY,
	user_id INTEGER REFERENCES Users(id) ON DELETE CASCADE, 
	emp_id INTEGER REFERENCES Employees(id) ON DELETE CASCADE,
	UNIQUE(user_id, emp_id)
)
`)
	if err != nil {
		return fmt.Errorf("failed to create subs table: %w", err)
	}
	log.Info("Tables created (or updated)")
	return nil

}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.Db.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.Db.Close()
}
