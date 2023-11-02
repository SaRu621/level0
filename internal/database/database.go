package database

import (
	"context"
	"fmt"

	"level0/internal/structs"

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

func CreateTableDelivery(db *pgx.Conn) error {
	_, err := db.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS delivery (
			id      INT PRIMARY KEY,
			name    VARCHAR(60),
			phone   VARCHAR(12),
			zip     VARCHAR(10),
			city    VARCHAR(20),
			address VARCHAR(30),
			region  VARCHAR(25),
			email   VARCHAR(50),
			FOREIGN KEY (id) REFERENCES info(id)
			)`)

	return err
}

func CreateTablePayment(db *pgx.Conn) error {
	_, err := db.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS payment (
			id 			  INT PRIMARY KEY, 
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
			FOREIGN KEY (id) REFERENCES info(id)
		)`)

	return err
}

func CreateTableItems(db *pgx.Conn) error {
	_, err := db.Exec(context.Background(),
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
			FOREIGN KEY (id) REFERENCES info(id)			
		)`)

	return err
}

func CreateTableInfo(db *pgx.Conn) error {
	_, err := db.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS info (
			id        		   SERIAL PRIMARY KEY,
			order_uid 		   VARCHAR(20),
			track_number 	   VARCHAR(14),
			entry 			   VARCHAR(10),
			locale 			   VARCHAR(2),
			internal_signature VARCHAR(20),
			customer_id 	   VARCHAR(20),
			delivery_service   VARCHAR(20),
			shardkey 		   VARCHAR(10),
			sm_id			   INT,
			date_created	   TIMESTAMP,
			oof_shard		   VARCHAR(10)
		)`)

	return err
}

func Migration(db *pgx.Conn) error {
	err1 := CreateTableInfo(db)
	err2 := CreateTableDelivery(db)
	err3 := CreateTablePayment(db)
	err4 := CreateTableItems(db)

	if err1 != nil {
		return err1
	}

	if err2 != nil {
		return err2
	}

	if err3 != nil {
		return err3
	}

	if err4 != nil {
		return err4
	}

	return nil
}

func WriteInDelivery(id int, data structs.Delivery, db *pgx.Conn) error {
	ToExec := fmt.Sprintf("INSERT INTO delivery VALUES (%d,'%s','%s','%s','%s','%s','%s','%s')",
		id, data.Name, data.Phone, data.Zip, data.City, data.Address, data.Region, data.Email)

	_, err := db.Exec(context.Background(), ToExec)

	return err
}

func WriteInPayment(id int, data structs.Payment, db *pgx.Conn) error {
	ToExec := fmt.Sprintf("INSERT INTO payment VALUES (%d, '%s','%s','%s','%s',%d,%d,'%s',%d,%d,%d)",
		id, data.Transaction, data.Request_id, data.Currency, data.Provider, data.Amount,
		data.Payment_dt, data.Bank, data.Delivery_cost, data.Goods_total, data.Custom_fee)

	_, err := db.Exec(context.Background(), ToExec)

	return err
}

func WriteInItems(id int, data structs.Items, db *pgx.Conn) error {
	ToExec := fmt.Sprintf("INSERT INTO items VALUES (%d, %d,'%s',%d,'%s','%s',%d,'%s',%d,%d,'%s',%d)",
		id, data.Chrt_id, data.Track_number, data.Price, data.Rid, data.Name,
		data.Sale, data.Size, data.Total_price, data.Nm_id, data.Brand, data.Status)

	_, err := db.Exec(context.Background(), ToExec)

	return err
}

func WriteInDatabase(data structs.Model, db *pgx.Conn) error {
	res, err := db.Query(context.Background(), "SELECT COUNT(*) FROM info")

	if err != nil {
		return err
	}

	var id int
	for res.Next() {
		res.Scan(&id)
	}
	id++

	ToExec := fmt.Sprintf("INSERT INTO info VALUES (%d, '%s','%s','%s','%s','%s','%s','%s','%s',%d,'%s','%s')",
		id, data.Order_uid, data.Track_number, data.Entry, data.Locale, data.Internal_signature,
		data.Customer_id, data.Delivery_service, data.ShardKey, data.Sm_id, data.Date_created,
		data.Oof_shard)

	_, err = db.Exec(context.Background(), ToExec)

	if err != nil {
		return err
	}

	WriteInDelivery(id, data.Deliv, db)
	WriteInPayment(id, data.Paym, db)

	for _, val := range data.Itms {
		WriteInItems(id, val, db)
	}

	return nil
}
