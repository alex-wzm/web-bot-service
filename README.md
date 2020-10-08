# web-bot-service

A super simple chatbot service implemented as a REST API.

Built with:
* [`go`](https://golang.org/)
* [`github.com/gorilla/mux`](https://github.com/gorilla/mux)
* [`cloud.google.com/go/dialogflow/apiv2`](https://godoc.org/cloud.google.com/go/dialogflow/apiv2)

It abstracts the gcloud Dialogflow API allowing easy integration with [web-bot-client](https://github.com/alex-wzm/web-bot-client) without exposing secret credentials which is an undersirable side effect of accessing Dialogflow's REST API directly within a client application.

## Project setup

### Generate and download Google Cloud credentials

Authention for the Dialogflow API uses service account key files decalred in your development environment as explained in [these docs](https://cloud.google.com/dialogflow/es/docs/quick/setup#auth)

### Compile and serve locally

```
go run main.go
```