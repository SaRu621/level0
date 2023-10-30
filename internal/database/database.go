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

func CreateTableDelivery(db *pgx.Conn) {
	db.Query(context.Background(),
		`CREATE TABLE IF NOT EXISTS delivery (
			id      INT,
			name    VARCHAR(60),
			phone   VARCHAR(12),
			zip     VARCHAR(10),
			city    VARCHAR(20),
			address VARCHAR(30),
			region  VARCHAR(25),
			email   VARCHAR(50),
			CONSTRAINTS pk_delivery_id PRIMARY KEY (id)
		)`)
}

func CreateTablePayment(db *pgx.Conn) {
	db.Query(context.Background(),
		`CREATE TABLE IF NOT EXISTS payment (
			id 			  INT,
			transaction   VARCHAR(20),
			request_id 	  VARCHAR(20),
			currency 	  VARCHAR(3),
			provider 	  VARCHAR(10),
			amout    	  INT NOT NULL,
			payment_dt 	  INT NOT NULL,
			bank 		  VARCHAR(10),
			delivery_cost INT NOT NULL,
			goods_total   INT NOT NULL,
			custom_fee    INT,
			CONSTRAINTS pk_payment_id PRIMARY KEY (id)
		)`)
}

func CreateTableItems(db *pgx.Conn) {
	db.Query(context.Background(),
		`CREATE TABLE IF NOT EXISTS items (
			id 			 INT,
			chrt_id 	 INT NOT NULL,
			track_number VARCHAR(14),
			price 		 INT NOT NULL,
			rid 		 VARCHAR(21),
			name 		 VARCHAR(50),
			sale 		 INT,
			size 		 VARCHAR(5),
			total_price  INT,
			nm_id 		 INT,
			brand 		 VARCHAR(50),
			status 		 INT,
			CONSTRAINTS pk_items_id PRIMARY KEY (id)
		)`)
}

func CreateTableInfo(db *pgx.Conn) {
	db.Query(context.Background(),
		`CREATE TABLE info (
			id INT,
			order_uid VARCHAR(20),
			track_number VARCHAR(14),
			entry VARCHAR(10),
			locale VARCHAR(2),
			internal_signature VARCHAR(20),
			customer_id VARCHAR(20),
			deli
			CONSTRAINTS pk_info_id PRIMARY KEY (id)
		)`)
}

func Migration(db *pgx.Conn) {
	CreateTableDelivery(db)
	CreateTableInfo(db)
	CreateTablePayment(db)
	CreateTableItems(db)
}
