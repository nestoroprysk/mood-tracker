package env

import "github.com/nestoroprysk/mood-tracker/internal/telegramclient"

// T is a bucket with configurations.
type T struct {
	telegramclient.Config `validate:"required"`
	Bucket                string `validate:"required"`
}
