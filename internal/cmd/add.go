package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"cloud.google.com/go/storage"

	"github.com/nestoroprysk/mood-tracker/internal/registry"
)

func newAdd(e env) (Cmd, error) {
	return func() (string, error) {
		j := userIDJSON(e.userID)

		b, err := e.Read(j)
		if errors.Is(err, storage.ErrObjectNotExist) {
			b = nil
		} else if err != nil {
			return "", err
		}

		r, err := registry.Make(b, e.userID)
		if err != nil {
			return "", err
		}

		r, err = r.WithItem(e.args...)
		if err != nil {
			return "", err
		}

		b, err = r.Dump()
		if err != nil {
			return "", err
		}

		if err := e.Override(j, b); err != nil {
			return "", err
		}

		return fmt.Sprintf("You have added %d entries that far. Good job!", len(r.Items)), nil
	}, nil
}

func userIDJSON(id int) string {
	return strconv.Itoa(id) + ".json"
}
