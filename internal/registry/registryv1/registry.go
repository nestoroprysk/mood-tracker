package registryv1

import (
	"time"

	"github.com/nestoroprysk/mood-tracker/internal/validator"
)

const Version = "v1"

// T is a registry with Version equal to v1.
type T struct {
	Version string `json:"version"`
	UserID  int    `json:"userid"`
	Items   []Item `json:"items"`
}

type Item struct {
	Time   time.Time `json:"time" validate:"required"`
	Mood   int       `json:"mood" validate:"min=1,max=5"`
	Labels []string  `json:"labels" validate:"dive,min=3"`
}

// FilterItems returns only valid items.
func FilterItems(is []Item) []Item {
	v := validator.New()

	var result []Item

	for _, i := range is {
		if err := v.Struct(i); err == nil {
			result = append(result, i)
		}
	}

	return result
}
