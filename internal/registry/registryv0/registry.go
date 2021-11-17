package registryv0

import "time"

// T is the zero unversioned registry.
// An empty file satisfies this version of registry.
type T []Item

type Item struct {
	Time   time.Time `json:"time"`
	Mood   int       `json:"mood"`
	Labels []string  `json:"labels"`
}
