Deploy Google-Cloud-Function:

gcloud functions deploy Firestore_HTTP_Add --entry-point=Add --runtime go111 --trigger-http --allow-unauthenticated --set-env-vars=PROJECT_ID=<project_id> --set-env-vars=PUBSUB_TOPIC=<topic>

NOTES...
Create local server: go build -o cmd/server cmd/main.go
https://codelabs.developers.google.com/codelabs/cloud-functions-go-http/#5