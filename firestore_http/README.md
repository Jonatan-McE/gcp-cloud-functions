Deploy Google-Cloud-Function:

gcloud functions deploy Firestore_ADD --entry-point=Add --runtime go111 --trigger-http --allow-unauthenticated --set-env-vars=GCLOUD_DATABASE_COLLECTION="<Database Collection>"
gcloud functions deploy Firestore_RETRIEVE --entry-point=Retrieve --runtime go111 --trigger-http --allow-unauthenticated --set-env-vars=GCLOUD_DATABASE_COLLECTION="<Database Collection>"

NOTES...
Create local server: go build -o cmd/server cmd/main.go
https://codelabs.developers.google.com/codelabs/cloud-functions-go-http/#5