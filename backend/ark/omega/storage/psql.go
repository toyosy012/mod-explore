package storage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

func ConnectPostgres(dsn string) (db *sqlx.DB, err error) {
	ctx := context.Background()
	{
		timeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		db, err = sqlx.ConnectContext(timeout, "postgres", dsn)
		if err != nil {
			return nil, err
		}
	}
	defer func() {
		if err != nil {
			db.Close()
		}
	}()

	{
		timeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = db.PingContext(timeout); err != nil {
			return nil, fmt.Errorf("connection timeout: %s ", err)
		}
	}

	return db, nil
}
