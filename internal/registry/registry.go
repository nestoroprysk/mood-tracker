package registry

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/getlantern/deepcopy"

	"github.com/nestoroprysk/mood-tracker/internal/registry/registryv0"
	"github.com/nestoroprysk/mood-tracker/internal/registry/registryv1"
	"github.com/nestoroprysk/mood-tracker/internal/validator"
)

// T is the current version of the registry.
//
// Registry is a mood diary for a person.
type T registryv1.T

// Make creates a new registry of the latest version.
//
// If the input body is empty, an empty registry is created.
func Make(b []byte, userID int) (T, error) {
	var regv1 registryv1.T
	if err := json.Unmarshal(b, &regv1); err == nil {
		if regv1.Version == registryv1.Version {
			regv1.Items = registryv1.FilterItems(regv1.Items)
			return T(regv1), nil
		}
	}

	var regv0 registryv0.T
	if len(b) != 0 {
		if err := json.Unmarshal(b, &regv0); err != nil {
			return T{}, err
		}
	}

	return T{
		Version: registryv1.Version,
		UserID:  userID,
		Items:   registryv1.FilterItems(regv0ItemsToRegv1Items(regv0)),
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

	i := registryv1.Item{
		Time:   time.Now().UTC(),
		Mood:   mood,
		Labels: args[1:],
	}

	v := validator.New()
	if err := v.Struct(i); err != nil {
		return T{}, err
	}

	result.Items = append(result.Items, i)

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
