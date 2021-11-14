package cmd

import (
	"fmt"
	"strings"

	"github.com/nestoroprysk/mood-tracker/internal/repository"
	"github.com/nestoroprysk/mood-tracker/internal/telegramclient"
)

// Config configures the command.
type Config struct {
	telegramclient.TelegramClient
	repository.Repository
	Args   string `validate:"required"`
	UserID int    `validate:"required"`
}

type config struct {
	telegramclient.TelegramClient
	repository.Repository
	args   []string `validate:"required"`
	userID int      `validate:"required"`
}

// Cmd executes the comment and returns the text result or error.
type Cmd func() (string, error)

// cmdCreator creates a command.
type cmdCreator func(config) (Cmd, error)

// registry maps command names to commands.
var registry = map[string]cmdCreator{
	"/add":  newAdd,
	"/stat": newStat,
}

func New(c Config) (Cmd, error) {
	tokens := parseTokens(c.Args)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("command name should be specified (indicate one of %s)", strings.Join(names(registry), " "))
	}

	creator, ok := registry[tokens[0]]
	if !ok {
		return nil, fmt.Errorf("command %s not found (indicate one of %s)", tokens[0], strings.Join(names(registry), " "))
	}

	cmd, err := creator(config{
		TelegramClient: c.TelegramClient,
		Repository:     c.Repository,
		args:           tokens[1:],
		userID:         c.UserID,
	})
	if err != nil {
		return nil, err
	}

	return cmd, nil
}

func names(r map[string]cmdCreator) []string {
	var result []string
	for n := range r {
		result = append(result, n)
	}

	return result
}

func parseTokens(m string) []string {
	ts := strings.Split(m, " ")
	var result []string
	for _, t := range ts {
		s := strings.TrimSpace(t)
		if s != "" {
			result = append(result, s)
		}
	}

	return result
}
