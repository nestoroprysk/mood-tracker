package cmd

import (
	"fmt"
	"strings"
)

// Cmd executes the comment and returns the text result or error.
type Cmd func() (string, error)

// cmdCreator creates a command.
type cmdCreator func(args ...string) (Cmd, error)

// registry maps command names to commands.
var registry = map[string]cmdCreator{
	"/add":  newAdd,
	"/stat": newStat,
}

func New(c string) (Cmd, error) {
	tokens := parseTokens(c)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("command name should be specified (indicate one of %s)", strings.Join(names(registry), " "))
	}

	creator, ok := registry[tokens[0]]
	if !ok {
		return nil, fmt.Errorf("command %s not found  (indicate one of %s)", tokens[0], strings.Join(names(registry), " "))
	}

	cmd, err := creator(tokens[1:]...)
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
