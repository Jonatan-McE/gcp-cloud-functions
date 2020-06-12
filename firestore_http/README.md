Deploy Google-Cloud-Function:

gcloud functions deploy \<function name\> --region=europe-west1 --allow-unauthenticated --entry-point=Submit --runtime go111 --set-env-vars=GCLOUD_DATABASE_COLLECTION="\<Database Collection\>" --set-env-vars=REDIRECT_URL="\<http handing page\>" --max-instances=3 --trigger-http 

gcloud functions deploy \<function name\> --region=europe-west1 --allow-unauthenticated --entry-point=Retrieve --runtime go111 --set-env-vars=GCLOUD_DATABASE_COLLECTION="\<Database Collection\>" --max-instances=3 --trigger-http

NOTES...
Create local server: go build -o cmd/server cmd/main.go
https://codelabs.developers.google.com/codelabs/cloud-functions-go-http/#5