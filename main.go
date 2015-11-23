package main

import (
	"log"
	"net/http"

	"github.com/karolgorecki/nbp-api/server"
)

func main() {

	rt := server.RegisterHandlers()
	log.Fatal(http.ListenAndServe(":8080", rt))
}
