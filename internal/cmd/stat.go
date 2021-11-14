package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/nestoroprysk/mood-tracker/internal/repository"
)

func newStat(r repository.Repository, userID int, args ...string) (Cmd, error) {
	return func() (string, error) {
		b, err := r.Read(userIDJSON(userID))
		if err != nil {
			return "", err
		}

		var reg Registry
		if err := json.Unmarshal(b, &reg); err != nil {
			return "", err
		}

		b, err = json.MarshalIndent(reg, "", " ")
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s", b), nil
	}, nil
}
