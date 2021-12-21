package main

import (
	"log"
	"net/http"

	sw "github.com/mbezuidenhout/tasmota.openapi/go"
)

func main() {

	router := sw.NewRouter()

	log.Fatal(http.ListenAndServe(":8080", router))
}
