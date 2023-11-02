package main

import (
	"encoding/json"
	"fmt"
	"os"

	"level0/internal/database"
	"level0/internal/structs"

	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
)

func main() {
	err := godotenv.Load("/home/rus/GoFolder/level0/.env")

	if err != nil {
		fmt.Println(err)
	}

	db, err := database.CreateConnection(os.Getenv("dbUser"), os.Getenv("dbPass"), os.Getenv("dbHost"), os.Getenv("dbPort"), os.Getenv("dbName"))

	if err != nil {
		fmt.Println(err)
	}

	err = database.Migration(db)

	if err != nil {
		fmt.Println(err)
	}

	sc, err := stan.Connect("test-cluster", "client-1") //подключение к кластеру NATS Streaming

	defer sc.Close()

	if err != nil {
		fmt.Println(err)
	}

	subscribtion, err := sc.Subscribe("my-channel", func(msg *stan.Msg) {
		var data structs.Model

		json.Unmarshal(msg.Data, &data)

		database.WriteInDatabase(data, db)
	}, stan.DurableName("my-durable"))

	defer subscribtion.Unsubscribe()

	if err != nil {
		fmt.Println(err)
	}

	ch := make(chan int)
	<-ch
}
