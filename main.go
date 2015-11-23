package main

import (
	"log"
	"net/http"
	"os"

	"github.com/karolgorecki/nbp-api/server"
)

func main() {

	rt := server.RegisterHandlers()
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), rt))
}
