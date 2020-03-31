// PROJECT_ID and PUBSUB_TOPIC need to be set as enviroment variables
// For ex: --set-env-vars=PROJECT_ID=littleglitch-lab --set-env-vars=PUBSUB_TOPIC=Test1

package pubsubhttp

import (
	"context"
	"net/http"
	"os"
	"regexp"

	"cloud.google.com/go/pubsub"
	"github.com/gorilla/schema"
)

// Get is an HTTP Cloud Function with a request parameter.
func Get(w http.ResponseWriter, r *http.Request) {
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

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	topic := client.Topic(os.Getenv("PUBSUB_TOPIC"))
	result := topic.Publish(ctx, &pubsub.Message{Data: []byte(reg.ReplaceAllString(q.ID, ""))})

	id, err := result.Get(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(id))
}
