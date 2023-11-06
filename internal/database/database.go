package database

import (
	"context"
	"fmt"
	"time"

	"level0/internal/structs"

	"github.com/jackc/pgx/v4"
)

func CreateConnection(user, pass, host, port, name string) (*pgx.Conn, error) { //создание подключения к БД
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

func CreateTableDelivery(db *pgx.Conn) error { //создание таблицы delivery
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

func CreateTablePayment(db *pgx.Conn) error { //создание таблицы payment
	_, err := db.Exec(context.Background(),
		`CREATE TABLE IF NOT EXISTS payment (
			id 			  INT PRIMARY KEY, 
			transaction   VARCHAR(20),
			request_id 	  VARCHAR(20),
			currency 	  VARCHAR(3),
			provider 	  VARCHAR(10),
			amount    	  INT NOT NULL,
			payment_dt 	  INT NOT NULL,
			bank 		  VARCHAR(10),
			delivery_cost INT NOT NULL,
			goods_total   INT NOT NULL,
			custom_fee    INT,
			FOREIGN KEY (id) REFERENCES info(id)
		)`)

	return err
}

func CreateTableItems(db *pgx.Conn) error { //создание таблицы items
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

func CreateTableInfo(db *pgx.Conn) error { //создание таблицы info
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

func Migration(db *pgx.Conn) error { //Миграция таблиц при их отсутствии
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

func WriteInDelivery(id int, data structs.Delivery, db *pgx.Conn) error { //запись данных в таблицу delivery
	ToExec := fmt.Sprintf("INSERT INTO delivery VALUES (%d,'%s','%s','%s','%s','%s','%s','%s')",
		id, data.Name, data.Phone, data.Zip, data.City, data.Address, data.Region, data.Email)

	_, err := db.Exec(context.Background(), ToExec)

	return err
}

func WriteInPayment(id int, data structs.Payment, db *pgx.Conn) error { //запись данных в таблицу payment
	ToExec := fmt.Sprintf("INSERT INTO payment VALUES (%d, '%s','%s','%s','%s',%d,%d,'%s',%d,%d,%d)",
		id, data.Transaction, data.Request_id, data.Currency, data.Provider, data.Amount,
		data.Payment_dt, data.Bank, data.Delivery_cost, data.Goods_total, data.Custom_fee)

	_, err := db.Exec(context.Background(), ToExec)

	return err
}

func WriteInItems(id int, data structs.Items, db *pgx.Conn) error { //запись данных в таблицу items
	ToExec := fmt.Sprintf("INSERT INTO items VALUES (%d, %d,'%s',%d,'%s','%s',%d,'%s',%d,%d,'%s',%d)",
		id, data.Chrt_id, data.Track_number, data.Price, data.Rid, data.Name,
		data.Sale, data.Size, data.Total_price, data.Nm_id, data.Brand, data.Status)

	_, err := db.Exec(context.Background(), ToExec)

	return err
}

func WriteInDatabase(data structs.Model, db *pgx.Conn) error { //запись данных в таблицу info и отсальные таблицы
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

func GetCash(cash map[string]structs.Model, db *pgx.Conn) error {
	rows, err := db.Query(context.Background(), `SELECT info.id, info.order_uid, info.track_number, info.entry, 
												info.locale, info.internal_signature, info.customer_id, 
												info.delivery_service, info.shardkey, info.sm_id, 
												info.date_created, info.oof_shard,
												delivery.name, delivery.phone, delivery.zip, 
												delivery.city, delivery.address, delivery.region, 
												delivery.email,
												payment.transaction, payment.request_id, payment.currency,
												payment.provider, payment.amount, payment.payment_dt, payment.bank, 
												payment.delivery_cost, payment.goods_total,payment.custom_fee,
												items.chrt_id, items.track_number, items.price, items.rid, items.name,
												items.sale, items.size, items.total_price, items.nm_id, items.brand, 
												items.status
												FROM info 
												INNER JOIN delivery ON info.id = delivery.id 
												INNER JOIN payment ON delivery.id = payment.id 
												INNER JOIN items ON payment.id = items.id`)

	if err != nil {
		return err
	}

	for rows.Next() {
		var id int
		var data structs.Model
		var chrt_id int
		var track_number string
		var price int
		var rid string
		var name string
		var sale int
		var size string
		var total_price int
		var nm_id int
		var brand string
		var status int
		var date_created time.Time

		err = rows.Scan(&id, &data.Order_uid, &data.Track_number, &data.Entry, &data.Locale,
			&data.Internal_signature, &data.Customer_id, &data.Delivery_service, &data.ShardKey, &data.Sm_id,
			&date_created, &data.Oof_shard, &data.Deliv.Name, &data.Deliv.Phone, &data.Deliv.Zip,
			&data.Deliv.City, &data.Deliv.Address, &data.Deliv.Region, &data.Deliv.Email,
			&data.Paym.Transaction, &data.Paym.Request_id, &data.Paym.Currency, &data.Paym.Provider,
			&data.Paym.Amount, &data.Paym.Payment_dt, &data.Paym.Bank, &data.Paym.Delivery_cost,
			&data.Paym.Goods_total, &data.Paym.Custom_fee,
			&chrt_id, &track_number, &price, &rid, &name, &sale, &size, &total_price, &nm_id,
			&brand, &status)

		if err != nil {
			fmt.Println(err)
			return err
		}

		data.Itms = append(data.Itms, structs.Items{
			Chrt_id:      chrt_id,
			Track_number: track_number,
			Price:        price,
			Rid:          rid,
			Name:         name,
			Sale:         sale,
			Size:         size,
			Total_price:  total_price,
			Nm_id:        nm_id,
			Brand:        brand,
			Status:       status,
		})

		data.Date_created = date_created.Format("02-01-2006 15:04:05")

		if _, ok := cash[data.Order_uid]; !ok {
			cash[data.Order_uid] = data

		} else {
			data = cash[data.Order_uid]
			data.Itms = append(cash[data.Order_uid].Itms, structs.Items{
				Chrt_id:      chrt_id,
				Track_number: track_number,
				Price:        price,
				Rid:          rid,
				Name:         name,
				Sale:         sale,
				Size:         size,
				Total_price:  total_price,
				Nm_id:        nm_id,
				Brand:        brand,
				Status:       status,
			})
			cash[data.Order_uid] = data
		}
	}

	return nil
}
