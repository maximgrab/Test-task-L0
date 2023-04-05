package main

import (
	"crypto/rand"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/stan.go"
)

type Order struct {
	//gorm.Model
	Uid  string //`gorm:"type:text;not null"` //`json:"uid" validate:"required,min=19,max=19"`
	Json string //`gorm:"type:text;not null"` //`json:json`
}

func main() {
	sc, err := stan.Connect(
		"test-cluster",
		"publisher",
		stan.NatsURL("nats://localhost:4222"),
	)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to NATS\n")
	for {
		time.Sleep(time.Second)
		bufUid := make([]byte, 19)
		bufJson := make([]byte, 100)
		rand.Read(bufUid)
		rand.Read(bufJson)
		for i := range bufUid {
			bufUid[i] = bufUid[i]%60 + '0'
		}
		for i := range bufJson {
			bufJson[i] = bufJson[i]%60 + '0'
		}
		var tmpOrder Order
		var buf []byte
		tmpOrder.Uid = string(bufUid)
		tmpOrder.Json = string(bufJson)
		buf, _ = json.Marshal(tmpOrder)
		if sc.Publish("orders", buf); err != nil {
			log.Fatal(err)
		}
		log.Printf("Order with uid:%v send to NATS chanel\n", string(bufUid))
	}

}
