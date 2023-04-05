package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/stan.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Order struct {
	//gorm.Model
	Uid  string `gorm:"type:text;not null" validate:"required,min=19,max=19"` //`json:"uid" validate:"required,min=19,max=19"`
	Json string `gorm:"type:text;not null"`                                   //`json:json`
}

func main() {
	dsn := "host=0.0.0.0 port=5432 user=postgres password=password dbname=orders sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()
	db.AutoMigrate(&Order{})
	var buf []Order

	ordersMap := make(map[string]string)
	db.Find(&buf)
	for _, i := range buf {
		ordersMap[i.Uid] = i.Json
		fmt.Println(ordersMap[i.Uid])
	}

	sc, err := stan.Connect(

		"test-cluster",
		"subscriber",
		stan.NatsURL("nats://localhost:4222"),
	)

	if err != nil {
		log.Fatal(err)
	}

	if sc.Publish("orders", []byte("All is Well")); err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		sc.Subscribe("orders", func(msg *stan.Msg) {
			var tmpOrder Order
			json.Unmarshal(msg.Data, &tmpOrder)
			ordersMap[tmpOrder.Uid] = tmpOrder.Json
			db.Create(&Order{Uid: tmpOrder.Uid, Json: tmpOrder.Json})
			fmt.Println("Uid: ", tmpOrder.Uid)
		})
	}()

	defer func(sc stan.Conn) {
		err = sc.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(sc)

	wg.Wait()
}
