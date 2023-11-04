package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"level0/internal/database"
	"level0/internal/structs"

	"github.com/gorilla/mux"
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

	cash := make(map[string]structs.Model)
	err = database.GetCash(cash, db)

	//fmt.Println(cash["9J0S4V5SSV0LQ1OZ929Y"])

	if err != nil {
		fmt.Println(err)
	}
	time.Sleep(time.Minute)

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

	r := mux.NewRouter()

	r.HandleFunc("/order/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userId := vars["id"]

		data := cash[userId]

		jsonData, err := json.Marshal(data)

		if err != nil {
			fmt.Println(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)

	}).Methods("GET")

	http.Handle("/", r)

	http.ListenAndServe(":"+os.Getenv("httpPort"), nil)
}
