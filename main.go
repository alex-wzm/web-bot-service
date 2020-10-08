package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"web-bot-service/detect_intent"

	"github.com/gorilla/mux"
)

var projectID string

// getProjectID extracts the projectID from the servive account JSON file
func getProjectID() string {
	data, err := ioutil.ReadFile("./credentials/service_account.json")
	if err != nil {
		fmt.Print(err)
	}

	type ServiceAccountJSONFile struct {
		FileType                string `json:"type"`
		ProjectID               string `json:"project_id"`
		PrivateKeyID            string `json:"private_key_id"`
		PrivateKey              string `json:"private_key"`
		ClientEmail             string `json:"client_email"`
		ClientID                string `json:"client_id"`
		AuthURI                 string `json:"auth_uri"`
		TokenURI                string `json:"token_uri"`
		AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
		ClientX509CertURL       string `json:"client_x509_cert_url"`
	}

	var obj ServiceAccountJSONFile

	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
	}

	fmt.Printf("Project ID: %s\n", obj.ProjectID)
	return obj.ProjectID

}

// detectIntent gets a response to a text query using the Dialogflow detectIntent API
func detectIntentQuery(query string) string {
	fmt.Printf("Query: \"%s\"\n", query)
	sessionID, languageCode := "123456789", "en"
	// TODO: generate new sessionID on start
	// TODO: accept languageCode as request parameter
	response, err := detect_intent.DetectIntentText(projectID, sessionID, query, languageCode)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Response: %s\n", response)
	return response
}

// rootHandler handles requests to the root route "/"
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch r.Method {
	case "GET":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "GET called"}`))
	case "POST":
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"message": "POST called"}`))
	case "PUT":
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(`{"message": "PUT called"}`))
	case "DELETE":
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "DELETE called"}`))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "not found"}`))
	}
}

// detectIntentHandler handles requests to the "/detect_intent" route
func detectIntentHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := detectIntentQuery("are you sentient?")
	w.Write([]byte(response))
}

func main() {
	fmt.Println("Starting web-bot-server...")
	projectID = getProjectID()

	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/", rootHandler)
	api.HandleFunc("/detect_intent", detectIntentHandler).Methods("POST")
	fmt.Println("\nweb-bot-server is ready")
	log.Fatal(http.ListenAndServe(":8081", r))
}
