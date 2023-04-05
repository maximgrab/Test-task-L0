package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/nats-io/stan.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Order struct {
	//gorm.Model
	Uid  string `gorm:"type:text;not null" validate:"required,min=19,max=19"`
	Json string `gorm:"type:text;not null"`
}
type OrderJson struct {
}

func main() {
	dsn := "host=0.0.0.0 port=5432 user=postgres password=password dbname=orders sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to database")
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
		log.Printf("Order with Uid:%v loaded from database\n", i.Uid)
	}
	log.Println("Data from database stored in local cache")
	sc, err := stan.Connect(

		"test-cluster",
		"subscriber",
		stan.NatsURL("nats://localhost:4222"),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to NAT")
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		sc.Subscribe("orders", func(msg *stan.Msg) {
			var tmpOrder Order
			err := json.Unmarshal(msg.Data, &tmpOrder)
			if err != nil {
				log.Fatal(err)
			}
			ordersMap[tmpOrder.Uid] = tmpOrder.Json
			log.Printf("Order with Uid:%v stored to local cache\n", tmpOrder.Uid)
			db.Create(&Order{Uid: tmpOrder.Uid, Json: tmpOrder.Json})
			log.Printf("Order with Uid:%v stored to database\n", tmpOrder.Uid)
		})
	}()
	log.Println("Subscribed to NATS")
	defer func(sc stan.Conn) {
		err = sc.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(sc)

	wg.Wait()
}
