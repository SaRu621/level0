package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"level0/internal/database"
	"level0/internal/structs"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nats-io/stan.go"
)

func main() {
	err := godotenv.Load("/home/rus/GoFolder/level0/.env") //файл переменных окружения

	if err != nil {
		fmt.Println(err)
	}

	db, err := database.CreateConnection(os.Getenv("dbUser"), os.Getenv("dbPass"), os.Getenv("dbHost"), os.Getenv("dbPort"), os.Getenv("dbName")) //Создание подключения к БД

	if err != nil {
		fmt.Println(err)
	}

	err = database.Migration(db) //Миграция базы данных при отсутствии соответствующих зависимостей

	if err != nil {
		fmt.Println(err)
	}

	cash := make(map[string]structs.Model) //кеш хранится в map
	err = database.GetCash(cash, db)       //подгружение кеша из бд

	if err != nil {
		fmt.Println(err)
	}

	sc, err := stan.Connect("test-cluster", "client-1") //подключение к кластеру NATS Streaming

	defer sc.Close()

	if err != nil {
		fmt.Println(err)
	}

	subscribtion, err := sc.Subscribe("my-channel", func(msg *stan.Msg) { //Подписка на канал "my-channel" NATS
		var data structs.Model

		err := json.Unmarshal(msg.Data, &data) //Валидация json-файла, при приеме невалидного json функция завершается

		if err != nil {
			return
		}

		database.WriteInDatabase(data, db) //Запись данных в БД
	}, stan.DurableName("my-durable"))

	defer subscribtion.Unsubscribe() //Отписка от канала перед закрытием приложения

	if err != nil {
		fmt.Println(err)
	}

	r := mux.NewRouter()

	r.HandleFunc("/order/{id}", func(w http.ResponseWriter, r *http.Request) { //Описание "ручки" (GET-запроса)
		vars := mux.Vars(r)
		userId := vars["id"]

		data := cash[userId]

		jsonData, err := json.Marshal(data) //Преорбразование данных из кеша в json

		if err != nil {
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData) //Отправка данных

	}).Methods("GET")

	http.Handle("/", r)

	http.ListenAndServe(":"+os.Getenv("httpPort"), nil)
}
