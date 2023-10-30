package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"
)

func CreateConnection(user, pass, host, port, name string) (*pgx.Conn, error) {
	config := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", user, pass, host, port, name)

	connConfig, err := pgx.ParseConfig(config)

	if err != nil {
		return nil, err
	}

	db, err := pgx.ConnectConfig(context.Background(), connConfig)

	if err != nil {
		return nil, err
	}

	return db, err
}

func CreateTable(db *pgx.Conn) {
	db.Query(context.Background(), `CREATE TABLE t`)
}

func Migration(db *pgx.Conn) {

}
