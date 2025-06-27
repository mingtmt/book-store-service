package initialize

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/mingtmt/book-store/pkg/logger"
)

var DB *pgxpool.Pool

func InitPostgres() (*pgxpool.Pool, error) {
	dsn := os.Getenv("DB_URL")

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		logger.Error("Unable to parse config", err, nil)
		os.Exit(1)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		logger.Error("Unable to connect to DB", err, nil)
		os.Exit(1)
	}

	if err = pool.Ping(context.Background()); err != nil {
		logger.Error("Unable to ping DB", err, nil)
		os.Exit(1)
	}

	DB = pool
	logger.Info("Postgres connection established successfully", nil)
	return DB, nil
}
