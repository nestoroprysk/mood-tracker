package registry

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/nestoroprysk/mood-tracker/internal/registry/registryv0"
	"github.com/nestoroprysk/mood-tracker/internal/registry/registryv1"

	"github.com/getlantern/deepcopy"
)

// T is the current version of the registry.
//
// Registry is a mood diary for a person.
type T registryv1.T

func Make(b []byte, userID int) (T, error) {
	var regv1 registryv1.T
	if err := json.Unmarshal(b, &regv1); err == nil {
		if regv1.Version == registryv1.Version {
			return T(regv1), nil
		}
	}

	var regv0 registryv0.T
	if err := json.Unmarshal(b, &regv0); err != nil {
		return T{}, err
	}

	return T{
		Version: registryv1.Version,
		UserID:  userID,
		Items:   regv0ItemsToRegv1Items(regv0),
	}, nil
}

func (t T) WithItem(args ...string) (T, error) {
	if len(args) == 0 {
		return T{}, fmt.Errorf("indicate mood from 1 to 5 with optional tags, e.g., /add 5 fascinated happy")
	}

	var result T
	if err := deepcopy.Copy(&result, t); err != nil {
		return T{}, err
	}

	mood, err := strconv.Atoi(args[0])
	if err != nil {
		return T{}, err
	}

	result.Items = append(result.Items, registryv1.Item{
		Time:   time.Now().UTC(),
		Mood:   mood,
		Labels: args[1:],
	})

	return result, nil
}

func (t T) Dump() ([]byte, error) {
	return json.MarshalIndent(t, "", " ")
}

func regv0ItemsToRegv1Items(is []registryv0.Item) []registryv1.Item {
	var result []registryv1.Item
	for _, i := range is {
		result = append(result, registryv1.Item{
			Time:   i.Time,
			Mood:   i.Mood,
			Labels: i.Labels,
		})
	}

	return result
}