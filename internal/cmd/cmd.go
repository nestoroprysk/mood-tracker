package cmd

import (
	"fmt"
	"strings"

	"github.com/nestoroprysk/mood-tracker/internal/repository"
	"github.com/nestoroprysk/mood-tracker/internal/telegramclient"
)

type Env struct {
	telegramclient.TelegramClient
	repository.Repository
	Args   string `validate:"required"`
	UserID int    `validate:"required"`
}

type env struct {
	telegramclient.TelegramClient
	repository.Repository
	args   []string `validate:"required"`
	userID int      `validate:"required"`
}

// Cmd executes the command and returns the text result or error.
type Cmd func() (string, error)

// cmdCreator creates a command.
type cmdCreator func(env) (Cmd, error)

// cmdRegistry maps command names to commands.
var cmdRegistry = map[string]cmdCreator{
	"/add":  newAdd,
	"/stat": newStat,
}

func New(e Env) (Cmd, error) {
	tokens := parseTokens(e.Args)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("command name should be specified (indicate one of %s)", strings.Join(names(cmdRegistry), " "))
	}

	creator, ok := cmdRegistry[tokens[0]]
	if !ok {
		return nil, fmt.Errorf("command %s not found (indicate one of %s)", tokens[0], strings.Join(names(cmdRegistry), " "))
	}

	cmd, err := creator(env{
		TelegramClient: e.TelegramClient,
		Repository:     e.Repository,
		args:           tokens[1:],
		userID:         e.UserID,
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
