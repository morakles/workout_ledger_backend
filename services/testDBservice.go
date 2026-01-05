package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type TestDBService struct {
	DB *sql.DB
}

func (s TestDBService) SelectAllTestRows() error {
	if s.DB == nil {
		return fmt.Errorf("database not initialized")
	}
	qctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	rows, err := s.DB.QueryContext(qctx, `SELECT id, "text" FROM testtable`)
	if err != nil {
		return fmt.Errorf("select failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id  int
			txt string
		)
		if err := rows.Scan(&id, &txt); err != nil {
			return fmt.Errorf("scan failed: %w", err)
		}
		log.Printf("row: id=%d name=%s", id, txt)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %w", err)
	}

	return nil
}
