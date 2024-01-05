package pgx

import (
	"context"
	"fmt"
	"log"

	"github.com/Sohbetbackend/selfProject/config"
	"github.com/jackc/pgx/v4/pgxpool"
)

type PgxStore struct {
	pool *pgxpool.Pool
}

func (d PgxStore) Pool() *pgxpool.Pool {
	return d.pool
}

func (d *PgxStore) Close() {
	d.pool.Close()
}

func Init() *PgxStore {
	connStr := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable connect_timeout=5", config.Conf.DbUsername, config.Conf.DbDatabase, config.Conf.DbPassword, config.Conf.DbHost, config.Conf.DbPort)
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	return &PgxStore{pool: pool}
}
