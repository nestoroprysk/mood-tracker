package registryv1

import "time"

const Version = "v1"

// T is a registry with Version equal to v1.
type T struct {
	Version string `json:"version"`
	UserID  int    `json:"userid"`
	Items   []Item `json:"items"`
}

type Item struct {
	Time   time.Time `json:"time"`
	Mood   int       `json:"mood"`
	Labels []string  `json:"labels"`
}
