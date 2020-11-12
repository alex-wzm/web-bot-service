package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"web-bot-service/detect_intent"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

const applicationName = "web-bot-client"
const port = ":8081"

var projectID string

// getProjectID extracts the projectID from the servive account JSON file
func getProjectID() string {
	data, err := ioutil.ReadFile("./credentials/service_account.json")
	if err != nil {
		log.Panic(err)
	}

	var parsed map[string]interface{}

	err = json.Unmarshal(data, &parsed)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("* projectID: %s\n", parsed["project_id"])
	return parsed["project_id"].(string)
}

// detectIntent gets a response to a text query using the Dialogflow detectIntent API
func detectIntentQuery(sessionID string, queryText string, languageCode string) (string, error) {
	log.Println("├ detectIntentQuery(")
	log.Println("│    sessionID:", sessionID)
	log.Println("│    queryText:", queryText)
	log.Println("│    languageCode:", languageCode)
	log.Println("│ )")

	response, err := detect_intent.DetectIntentText(projectID, sessionID, queryText, languageCode)
	if err != nil {
		fmt.Println(err.Error())
		return "", errors.New("dialogflow error")
	}
	log.Printf("└ Response: \"%s\"\n", response)
	return response, nil
}

// rootHandler handles requests to the root route "/"
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: refactor to use as middleware for handler logging and writing headers
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
	switch r.Method {
	case http.MethodPost:
		log.Println("[POST] /detect_intent")
		type QueryBody struct {
			SessionID    string `json:",omitempty"`
			QueryText    string
			LanguageCode string
		}

		var body QueryBody
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			log.Print(err)
			return
		}
		log.Printf("├ Request body: %+v\n", body)

		// guards
		if body.QueryText == "" {
			// required
			log.Println("└ X Missing query text")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing query text."))
			return
		}
		if body.LanguageCode == "" {
			// required
			log.Println("└ X Missing language code")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Missing language code. Doc: https://dialogflow.com/docs/reference/language"))
			return
		}
		if body.SessionID == "" {
			// optional - generate single-use sessionID if none is provided
			log.Println("├ ✓ Generating default sessionID...")
			sessionID := fmt.Sprintf("default_session_%s", uuid.New())
			body.SessionID = sessionID
		}

		response, err := detectIntentQuery(body.SessionID, body.QueryText, body.LanguageCode)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not found."))
	}
}

func main() {
	log.Printf("* Starting %s...", applicationName)
	projectID = getProjectID()

	// initalize router + subrouter prefix
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()

	// handle routes
	api.HandleFunc("/", rootHandler)
	api.HandleFunc("/detect_intent", detectIntentHandler).Methods(http.MethodPost)

	// enable CORS
	// TODO: set custom options for whitelisting known applications only
	handler := cors.Default().Handler(api)

	log.Printf("* %s is running at port %s", applicationName, port)
	log.Fatal(http.ListenAndServe(port, handler))
}
