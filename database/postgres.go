package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func (cfg Config) connString() string {
	ssl := cfg.SSLMode
	if ssl == "" {
		ssl = "disable"
	}
	port := cfg.Port
	if port == "" {
		port = "5432"
	}

	// Keyword/value style avoids URL parsing issues with special chars.
	// We still set User/Password directly on the ConnConfig below.
	return fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s", cfg.Host, port, cfg.DBName, ssl)
}

// Open creates a database/sql DB backed by pgx and validates it with PingContext.
// It supports special characters in username/password because they are set on pgx.ConnConfig directly.
func Open(ctx context.Context, cfg Config) (*sql.DB, error) {
	baseConnStr := cfg.connString()

	connCfg, err := pgx.ParseConfig(baseConnStr)
	if err != nil {
		return nil, fmt.Errorf("parse pgx config: %w", err)
	}

	// Set credentials explicitly so special chars don't need escaping.
	connCfg.User = cfg.User
	connCfg.Password = cfg.Password

	registered := stdlib.RegisterConnConfig(connCfg)
	db, err := sql.Open("pgx", registered)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("connected to the database successfully")
	return db, nil
}
