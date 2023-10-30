package main

import (
    "fmt"
    "github.com/nats-io/stan.go"
)

func main() {
    sc, err := stan.Connect("test-cluster", "client-1") //подключение к кластеру NATS Streaming

    defer sc.Close()

    if err != nil {
        fmt.Println(err)
    }

    subscribtion, err := sc.Subscribe("my-channel", func(msg *stan.Msg) {
        fmt.Printf("Получено сообщение: %s\n", string(msg.Data))
    }, stan.DurableName("my-durable"))

    defer subscribtion.Unsubscribe()

    if err != nil {
        fmt.Println(err)
    }

    ch := make(chan int) 
    <- ch
}
