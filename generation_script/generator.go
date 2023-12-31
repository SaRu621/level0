package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"level0/internal/structs"

	"github.com/nats-io/stan.go"
)

func Generator() structs.Model {
	rand.Seed(time.Now().UnixNano())

	tn := StringGenerator(14)

	res := structs.Model{
		Order_uid:          StringGenerator(20),
		Track_number:       tn,
		Entry:              StringGenerator(10),
		Locale:             StringGenerator(2),
		Internal_signature: StringGenerator(20),
		Customer_id:        StringGenerator(20),
		Delivery_service:   StringGenerator(20),
		ShardKey:           StringGenerator(10),
		Sm_id:              NumberGenerator(8),
		Date_created:       time.Now(),
		Oof_shard:          StringGenerator(10),
		Paym: structs.Payment{
			Transaction:   StringGenerator(20),
			Request_id:    StringGenerator(20),
			Currency:      StringGenerator(3),
			Provider:      StringGenerator(10),
			Amount:        NumberGenerator(4),
			Payment_dt:    NumberGenerator(5),
			Bank:          StringGenerator(5),
			Delivery_cost: NumberGenerator(3),
			Goods_total:   NumberGenerator(2),
			Custom_fee:    NumberGenerator(3),
		},
		Deliv: structs.Delivery{
			Name:    StringGenerator(7) + " " + StringGenerator(7),
			Phone:   "+79251900718",
			Zip:     StringGenerator(10),
			City:    StringGenerator(12),
			Address: StringGenerator(20),
			Region:  StringGenerator(10),
			Email:   StringGenerator(10) + "@" + StringGenerator(5) + ".ru",
		},
	}

	count := rand.Intn(5) + 2

	res.Itms = make([]structs.Items, count)

	for i, _ := range res.Itms {
		res.Itms[i] = structs.Items{
			Chrt_id:      NumberGenerator(5),
			Track_number: tn,
			Price:        NumberGenerator(5),
			Rid:          StringGenerator(21),
			Name:         StringGenerator(20),
			Sale:         NumberGenerator(2),
			Size:         StringGenerator(1),
			Total_price:  NumberGenerator(4),
			Nm_id:        NumberGenerator(7),
			Brand:        StringGenerator(7),
			Status:       202,
		}
	}

	return res
}

func StringGenerator(n int) string {
	str := ""

	for i := 0; i < n; i++ {
		flag := rand.Intn(2)

		if flag == 0 {
			str += string(rune((65 + rand.Intn(26))))
		} else {
			str += strconv.Itoa(rand.Intn(10))
		}
	}

	return str
}

func NumberGenerator(n int) int {
	str := ""

	for i := 0; i < n; i++ {
		str += strconv.Itoa(rand.Intn(10))
	}

	num, _ := strconv.Atoi(str)

	return num
}

func Loader(sc *stan.Conn, wg *sync.WaitGroup) {
	inst := Generator()
	jsondata, err := json.Marshal(inst)

	if err != nil {
		fmt.Println(err)
	}

	(*sc).Publish("my-channel", jsondata)

	defer wg.Done()
}

const n = 10

func main() {
	clusterId := "test-cluster"
	clientId := "client-2"
	//channel := "my-channel"

	sc, err := stan.Connect(clusterId, clientId)

	if err != nil {
		fmt.Println(err)
	}

	defer sc.Close()

	for {
		var wg sync.WaitGroup
		wg.Add(n)
		for i := 0; i < n; i++ {
			go Loader(&sc, &wg)
		}
		wg.Wait()
		time.Sleep(time.Second * 10)
	}
}
