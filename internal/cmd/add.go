package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

// Registry is a list of items.
type Registry []Item

// Item is a mood item.
type Item struct {
	Time   time.Time `json:"time"`
	Mood   int       `json:"mood"`
	Labels []string  `json:"labels"`
}

func newAdd(c config) (Cmd, error) {
	return func() (string, error) {
		if len(c.args) == 0 {
			return "", errors.New("indicate the mood valued from 1 to 5 and optional tags, e.g., /add 5 happy energetic")
		}

		mood, err := strconv.Atoi(c.args[0])
		if err != nil {
			return "", errors.New("indicate the mood valued from 1 to 5 and optional tags, e.g., /add 5 happy energetic")
		}

		i := Item{
			Time:   time.Now().UTC(),
			Mood:   mood,
			Labels: c.args[1:],
		}

		b, err := c.Read(userIDJSON(c.userID))
		if err != nil {
			return "", err
		}

		var r Registry
		if err := json.Unmarshal(b, &r); err != nil {
			return "", err
		}
		r = append(r, i)

		b, err = json.MarshalIndent(r, "", " ")
		if err != nil {
			return "", err
		}

		if err := c.Override(userIDJSON(c.userID), b); err != nil {
			return "", err
		}

		return fmt.Sprintf("Added a new entry. Thank you for using the bot! You have added %d entries that far. Good job! Enter /stat to see analytics on your mood.", len(r)), nil
	}, nil
}

func userIDJSON(id int) string {
	return strconv.Itoa(id) + ".json"
}
