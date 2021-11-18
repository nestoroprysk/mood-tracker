package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/nestoroprysk/mood-tracker/internal/cmd"
	"github.com/nestoroprysk/mood-tracker/internal/env"
	"github.com/nestoroprysk/mood-tracker/internal/repository"
	"github.com/nestoroprysk/mood-tracker/internal/telegramclient"
)

func MoodTracker(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			respond(w, nil, fmt.Sprintf("%s", r))
		}
	}()

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

	logUpdate(u)

	t := telegramclient.New(env.Config, u.Message.Chat.ID, http.DefaultClient)

	repo, err := repository.New(env.Bucket)
	if err != nil {
		respond(w, t, err.Error())
		return
	}

	c, err := cmd.New(cmd.Env{
		TelegramClient: t,
		Repository:     repo,
		UserID:         u.Message.From.ID,
		Args:           u.Message.Text,
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

func logUpdate(u telegramclient.Update) {
	res, err := json.MarshalIndent(u, "", "")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(strings.ReplaceAll(string(res), "\n", ""))
}

func respond(writer http.ResponseWriter, client telegramclient.TelegramClient, response string) {
	log.Println(strings.ReplaceAll(response, "\n", ""))

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Println(err)
	}

	if client != nil && response != "" {
		r, err := client.SendMessage(response)
		if err != nil {
			log.Println(err)
			return
		}

		log.Println(r)
	} else if response != "" {
		log.Println(response)
	} else {
		log.Println("empty response")
	}
}
