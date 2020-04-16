package firestorehttp

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type atendy struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

// Retrieve is an HTTP Cloud Function thar return database docuuments in json.
func Retrieve(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GCLOUD_PROJECT"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer client.Close()

	var atendys []atendy
	iter := client.Collection(os.Getenv("GCLOUD_DATABASE_COLLECTION")).OrderBy("timestamp", firestore.Desc).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		atendys = append(atendys, atendy{ID: doc.Data()["id"].(string), Timestamp: doc.Data()["timestamp"].(time.Time)})
	}
	jsonData, err := json.Marshal(atendys)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)

	return
}
