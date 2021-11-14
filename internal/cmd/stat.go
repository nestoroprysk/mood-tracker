package cmd

import (
	"fmt"

	"github.com/nestoroprysk/mood-tracker/internal/repository"
)

func newStat(r repository.Repository, userID int, args ...string) (Cmd, error) {
	return func() (string, error) {
		return fmt.Sprintf("%v", args), nil
	}, nil
}
