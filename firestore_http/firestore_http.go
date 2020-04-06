package firestorehttp

import (
	"context"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/firestore"
)

// Add is an HTTP Cloud Function with a URL parameter.
func Add(w http.ResponseWriter, r *http.Request) {

	filter, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	keys, ok := r.URL.Query()["id"]
	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	value := filter.ReplaceAllString(keys[0], "")

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GCLOUD_PROJECT"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer client.Close()

	_, _, err = client.Collection(os.Getenv("GCLOUD_DATABASE_COLLECTION")).Add(ctx, map[string]interface{}{
		"timestamp": firestore.ServerTimestamp,
		"id":        value,
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Thank you for your submission\n"))

	return
}
