package main

import (
	"fmt"
	"log"
	"net/http"

	firestorehttp "github.com/jonatan-mce/gcp-cloud-functions/firestore_http"
)

func main() {
	http.HandleFunc("/", firestorehttp.Add)
	fmt.Println("Listening on localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
