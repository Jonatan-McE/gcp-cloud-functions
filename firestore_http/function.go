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
	UserID     string    `json:"userid"`
	OccasionID string    `json:"occasionid"`
	Timestamp  time.Time `json:"timestamp"`
}

// Submit new user info into the database
func Submit(w http.ResponseWriter, r *http.Request) {

	userID, okUserID := r.URL.Query()["userid"]
	occasionID, okoccasionID := r.URL.Query()["occasionid"]
	if okUserID || okoccasionID {
		var user string
		var occasion string

		if len(userID[0]) < 1 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Missing userID value\n"))
			return
		}
		if len(occasionID[0]) < 1 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Missing occasionID value\n"))
			return
		}
		filter, err := regexp.Compile("[^a-zA-Z0-9]+")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		user = filter.ReplaceAllString(userID[0], "")
		occasion = filter.ReplaceAllString(occasionID[0], "")

		ctx := context.Background()
		client, err := firestore.NewClient(ctx, os.Getenv("GCLOUD_PROJECT"))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		defer client.Close()

		_, _, err = client.Collection(os.Getenv("GCLOUD_DATABASE_COLLECTION")).Add(ctx, map[string]interface{}{
			"timestamp":  firestore.ServerTimestamp,
			"userid":     user,
			"occasionid": occasion,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		http.Redirect(w, r, os.Getenv("REDIRECT_URL"), 307)

		return
	}

	w.WriteHeader(http.StatusNotFound)

	return
}

// Retrieve user info from the database
func Retrieve(w http.ResponseWriter, r *http.Request) {

	var days int
	var atendys []atendy

	keys, ok := r.URL.Query()["days"]
	if ok {
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
	if ok == false {
		iter = client.Collection(os.Getenv("GCLOUD_DATABASE_COLLECTION")).Where("timestamp", ">", time.Now().AddDate(0, 0, -1)).OrderBy("timestamp", firestore.Desc).Documents(ctx)
	} else if days > 0 {
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

		atendys = append(atendys, atendy{UserID: doc.Data()["userid"].(string), OccasionID: doc.Data()["occasionid"].(string), Timestamp: doc.Data()["timestamp"].(time.Time)})
	}
	/*
		if atendys == nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[]"))
			return
		}
	*/
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
