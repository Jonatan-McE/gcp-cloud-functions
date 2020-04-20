Deploy Google-Cloud-Function:

gcloud functions deploy \<function name\> --region=europe-west1 --entry-point=Get --runtime go111 --trigger-http --allow-unauthenticated --set-env-vars=GCLOUD_DATABASE_COLLECTION="\<Database Collection\>" REDIRECT_URL="\<http handing page\>"

NOTES...
Create local server: go build -o cmd/server cmd/main.go
https://codelabs.developers.google.com/codelabs/cloud-functions-go-http/#5