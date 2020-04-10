package firestorehttp

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// Retrieve is an HTTP Cloud Function with a URL parameter.
func Retrieve(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	client, err := firestore.NewClient(ctx, os.Getenv("GCLOUD_PROJECT"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	defer client.Close()

	iter := client.Collection(os.Getenv("GCLOUD_DATABASE_COLLECTION")).Documents(ctx)
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
		fmt.Println(doc.Data())
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Looks good\n"))

	return
}
