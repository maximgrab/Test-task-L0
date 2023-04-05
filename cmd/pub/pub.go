package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	structs "github.com/maximgrab/Test-task-L0/pkg"
	"github.com/nats-io/stan.go"
)

func main() {
	jsonFile, err := os.Open("model.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var jsonModel structs.OrderJson
	json.Unmarshal(byteValue, &jsonModel)

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
		rand.Read(bufUid)
		for i := range bufUid {
			bufUid[i] = bufUid[i]%25 + 'a'
		}

		var tmpOrder structs.Order
		var buf []byte
		jsonModel.OrderUid = string(bufUid)
		tmpOrder.Uid = string(jsonModel.OrderUid)
		buf, _ = json.Marshal(jsonModel)
		tmpOrder.Json = string(buf)
		buf, _ = json.Marshal(tmpOrder)
		if sc.Publish("orders", buf); err != nil {
			log.Fatal(err)
		}
		log.Printf("Order with uid:%v send to NATS chanel\n", string(bufUid))
	}

}
