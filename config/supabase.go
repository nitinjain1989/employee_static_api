package config

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() *pgxpool.Pool {
	connStr := os.Getenv("SUPABASE_DIRECT_CONNECTION")

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		panic(err)
	}

	// 🔥 Disable prepared statement cache (IMPORTANT for dynamic queries)
	config.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		panic(err)
	}

	return db
}
