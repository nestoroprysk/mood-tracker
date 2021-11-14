package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/nestoroprysk/mood-tracker/internal/cmd"
	"github.com/nestoroprysk/mood-tracker/internal/env"
	"github.com/nestoroprysk/mood-tracker/internal/repository"
	"github.com/nestoroprysk/mood-tracker/internal/telegramclient"
)

func MoodTracker(w http.ResponseWriter, r *http.Request) {
	env := env.T{
		Config: telegramclient.Config{
			Token: os.Getenv("MOOD_TRACKER_BOT_TOKEN"),
		},
		Bucket: os.Getenv("MOOD_TRACKER_BUCKET"),
	}

	u, err := telegramclient.ParseUpdate(r.Body)
	if err != nil {
		respond(w, nil, err.Error())
		return
	}

	t := telegramclient.New(env.Config, u.Message.Chat.ID, http.DefaultClient)
	if err != nil {
		respond(w, nil, err.Error())
		return
	}

	repo, err := repository.New(env.Bucket)
	if err != nil {
		respond(w, t, err.Error())
		return
	}

	c, err := cmd.New(cmd.Config{
		Repository: repo,
		UserID:     u.Message.From.ID,
		Args:       u.Message.Text,
	})
	if err != nil {
		respond(w, t, err.Error())
		return
	}

	result, err := c()
	if err != nil {
		respond(w, t, err.Error())
		return
	}

	respond(w, t, result)
}

func respond(writer http.ResponseWriter, client telegramclient.TelegramClient, response string) {
	log.Println(response)

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Println(err)
	}

	r, err := client.Send(response)
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(r)
}
