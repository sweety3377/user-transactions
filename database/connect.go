package database

import (
	"blackwallgroup/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Connect(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	connString := getConnString(cfg)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, nil
}

func getConnString(cfg *config.Config) string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DbHost, cfg.DbPort, cfg.DbLogin, cfg.DbPassword, cfg.DbName)
}
