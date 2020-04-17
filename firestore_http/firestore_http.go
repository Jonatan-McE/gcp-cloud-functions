package firestorehttp

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

type atendy struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
}

// Get ...
func Get(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["id"]
	if ok {
		var id string

		if len(keys[0]) < 1 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Missing ID value\n"))
			return
		}
		filter, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		id = filter.ReplaceAllString(keys[0], "")

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
			"id":        id,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("OK"))

		return
	}

	keys, ok = r.URL.Query()["days"]
	if ok {
		var days int
		var atendys []atendy

		if len(keys[0]) < 1 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Missing days value\n"))
			return
		}

		filter, err := regexp.Compile("[^0-9]+")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		days, err = strconv.Atoi(filter.ReplaceAllString(keys[0], ""))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Missing days value\n"))
			w.Write([]byte(err.Error()))
			return
		}

		ctx := context.Background()
		client, err := firestore.NewClient(ctx, os.Getenv("GCLOUD_PROJECT"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer client.Close()

		var iter *firestore.DocumentIterator
		if days != 0 {
			iter = client.Collection(os.Getenv("GCLOUD_DATABASE_COLLECTION")).Where("timestamp", ">", time.Now().AddDate(0, 0, (days*-1))).OrderBy("timestamp", firestore.Desc).Documents(ctx)
		} else {
			iter = client.Collection(os.Getenv("GCLOUD_DATABASE_COLLECTION")).OrderBy("timestamp", firestore.Desc).Documents(ctx)
		}
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

	w.WriteHeader(http.StatusNotFound)

	return
}
