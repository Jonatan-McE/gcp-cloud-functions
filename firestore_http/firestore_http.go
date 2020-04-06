package firestorehttp

import (
	"context"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/firestore"
	"github.com/gorilla/schema"
)

// Add is an HTTP Cloud Function with a request parameter.
func Add(w http.ResponseWriter, r *http.Request) {
	var q struct {
		ID string `json:"id"`
	}

	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	var decoder = schema.NewDecoder()
	err = decoder.Decode(&q, r.URL.Query())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if q.ID != "" {
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
			"id":        reg.ReplaceAllString(q.ID, ""),
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("OK"))

	}

	return
}
