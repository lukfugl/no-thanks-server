package main

import (
	"log"
	"net/http"

	"./api"
)

func main() {
	http.HandleFunc("/", api.ExampleHandler)
	log.Println("** Service Started on Port 8080 **")
	err := http.ListenAndServeTLS(":8080", "certs/https-server.crt", "certs/https-server.key", nil)
	if err != nil {
		log.Fatal(err)
	}
}
