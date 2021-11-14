package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/nestoroprysk/mood-tracker/internal/repository"
)

// Registry is a list of items.
type Registry []Item

// Item is a mood item.
type Item struct {
	Time   time.Time `json:"time"`
	Mood   int       `json:"mood"`
	Labels []string  `json:"labels"`
}

func newAdd(r repository.Repository, userID int, args ...string) (Cmd, error) {
	return func() (string, error) {
		if len(args) == 0 {
			return "", errors.New("indicate the mood valued from 1 to 5 and optional tags, e.g., /add 5 happy energetic")
		}

		mood, err := strconv.Atoi(args[0])
		if err != nil {
			return "", errors.New("indicate the mood valued from 1 to 5 and optional tags, e.g., /add 5 happy energetic")
		}

		i := Item{
			Time:   time.Now().UTC(),
			Mood:   mood,
			Labels: args[1:],
		}

		b, err := r.Read(userIDJSON(userID))
		if err != nil {
			return "", err
		}

		var reg Registry
		if err := json.Unmarshal(b, &reg); err != nil {
			return "", err
		}
		reg = append(reg, i)

		b, err = json.MarshalIndent(reg, "", " ")
		if err != nil {
			return "", err
		}

		if err := r.Override(userIDJSON(userID), b); err != nil {
			return "", err
		}

		return fmt.Sprintf("Added a new entry. Thank you for using the bot! You have added %d entries that far. Good job! Enter /stat to see analytics on your mood.", len(reg)), nil
	}, nil
}

func userIDJSON(id int) string {
	return strconv.Itoa(id) + ".json"
}
