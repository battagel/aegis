package postgres

import (
	"aegis/pkg/logger"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type CloseFunc func()

type PostgresDB struct {
	logger logger.Logger
	pool   *pgxpool.Pool
}

func CreatePostgresDB(logger logger.Logger, user, password, endpoint string, database string) (*PostgresDB, CloseFunc, error) {
	logger.Debugw("Connecting to Postgresql DB",
		"database", database,
	)
	connectionUrl := fmt.Sprintf(
		"postgres://%v:%v@%v/%v?sslmode=disable",
		user, password, endpoint, database,
	)
	pool, err := pgxpool.New(context.Background(), connectionUrl)
	if err != nil {
		logger.Errorw("Error connecting to Postgres DB",
			"user", user,
			"password", password,
			"endpoint", endpoint,
			"database", database,
			"error", err,
		)
		return nil, nil, err
	}

	return &PostgresDB{
		logger: logger,
		pool:   pool,
	}, pool.Close, nil
}

func (p *PostgresDB) CreateTable(tableName string) error {
	query := fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			ID SERIAL PRIMARY KEY,
			ObjectKey TEXT NOT NULL,
			BucketName TEXT NOT NULL,
			Result TEXT NOT NULL,
			Antivirus TEXT NOT NULL,
			Timestamp TIMESTAMP NOT NULL
		);`, tableName)
	_, err := p.pool.Exec(context.Background(), query)
	if err != nil {
		p.logger.Errorw("Error creating table",
			"error", err,
		)
		return err
	}
	p.logger.Debugw("Table is created",
		"table", tableName,
	)
	return nil
}

func (p *PostgresDB) Insert(tableName, bucketName, objectKey, result, antivirus, timestamp string) error {
	query := fmt.Sprintf("INSERT INTO %s (ObjectKey, BucketName, Result, Antivirus, Timestamp) VALUES ($1, $2, $3, $4, $5)", tableName)
	_, err := p.pool.Exec(context.Background(), query, objectKey, bucketName, result, antivirus, timestamp)
	if err != nil {
		p.logger.Errorw("Error inserting into table",
			"error", err,
		)
		return err
	}
	return nil
}
