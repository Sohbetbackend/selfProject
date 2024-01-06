package pgx

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/Sohbetbackend/selfProject/config"
	"github.com/Sohbetbackend/selfProject/internal/models"
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

type pgxQuery func(conn *pgxpool.Conn) (err error)

func (d *PgxStore) runQuery(ctx context.Context, f pgxQuery) (err error) {
	err = d.Pool().AcquireFunc(ctx, f)
	if err != nil {
		return err
	}
	return
}

func parseColumnsForScan(sub interface{}, addColumns ...interface{}) []interface{} {
	s := reflect.ValueOf(sub).Elem()
	numCols := s.NumField() - len(sub.(models.HasRelationFields).RelationFields())

	columns := []interface{}{}
	for i := 0; i < numCols; i++ {
		field := s.Field(i)
		columns = append(columns, field.Addr().Interface())
	}
	columns = append(columns, addColumns...)
	return columns
}
