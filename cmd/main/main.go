package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	structs "github.com/maximgrab/Test-task-L0/pkg"
	"github.com/nats-io/stan.go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	db.AutoMigrate(&structs.Order{})
	var buf []structs.Order
	ordersMap := make(map[string]string)
	db.Find(&buf)
	for _, i := range buf {
		ordersMap[i.Uid] = i.Json
		log.Printf("Order with Uid:%v loaded from database\n", i.Uid)
	}
	log.Println("Data from database stored in local cache")
	//Тут сервачок пилим
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case "GET":
				http.ServeFile(w, r, "form.html")
			case "POST":
				uid := r.FormValue("uid")
				log.Printf("http server: order with uid %v\n requested", uid)
				fmt.Fprintf(w, "%v", ordersMap[uid])
			}
		})
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("http server started")
	}()
	sc, err := stan.Connect(

		"test-cluster",
		"subscriber",
		stan.NatsURL("nats://localhost:4222"),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to NAT")
	wg.Add(1)

	go func() {
		sc.Subscribe("orders", func(msg *stan.Msg) {
			var tmpOrder structs.Order
			err := json.Unmarshal(msg.Data, &tmpOrder)
			if err != nil {
				log.Fatal(err)
			}
			ordersMap[tmpOrder.Uid] = tmpOrder.Json
			log.Printf("Order with Uid:%v stored to local cache\n", tmpOrder.Uid)
			db.Create(&structs.Order{Uid: tmpOrder.Uid, Json: tmpOrder.Json})
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
