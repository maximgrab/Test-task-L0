package main

import (
	"log"
	"net/http"
)

func main() {
	srv := new(http.Server)
	if err := srv.Run("8080"); err != nil {
		log.Fatal("http server running error: %s", err.Error())
	}
}
