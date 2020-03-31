package main

import (
	"fmt"
	"log"
	"net/http"

	pubsubhttp "github.com/jonatan-mce/gcp-cloud-functions/pubsub_http"
)

func main() {
	http.HandleFunc("/", pubsubhttp.Get)
	fmt.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
