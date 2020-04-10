package main

import (
	"fmt"
	"log"
	"net/http"

	firestorehttp "github.com/Jonatan-McE/gcp-cloud-functions/firestore_retrieve_http"
)

func main() {
	http.HandleFunc("/", firestorehttp.Retrieve)
	fmt.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}