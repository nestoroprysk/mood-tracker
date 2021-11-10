package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nestoroprysk/mood-tracker/internal/env"
	"github.com/nestoroprysk/mood-tracker/internal/telegramclient"
)

func MoodTracker(w http.ResponseWriter, r *http.Request) {
	env := env.T{
		Config: telegramclient.Config{
			Token: os.Getenv("MOOD_TRACKER_BOT_TOKEN"),
		},
	}

	u, err := telegramclient.ParseUpdate(r.Body)
	if err != nil {
		respond(w, fmt.Sprintf("failed to parse the update: %s", err.Error()), http.StatusBadRequest)
		return
	}

	result, err := json.MarshalIndent(u, "", " ")
	if err != nil {
		respond(w, fmt.Sprintf("failed to parse the response: %s", err.Error()), http.StatusBadRequest)
		return
	}

	t := telegramclient.New(env.Config, u.Message.Chat.ID, http.DefaultClient)
	if err != nil {
		respond(w, fmt.Sprintf("failed to parse the update: %s", err.Error()), http.StatusBadRequest)
		return
	}

	response, err := t.Send(string(result))
	if err != nil {
		respond(w, fmt.Sprintf("failed to send the message to telegram: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	respond(w, response, http.StatusOK)
}

func respond(writer http.ResponseWriter, response interface{}, code int) error {
	log.Printf("%s %d", response, code)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	return json.NewEncoder(writer).Encode(response)
}
