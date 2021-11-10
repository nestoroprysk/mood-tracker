package env

import "github.com/nestoroprysk/mood-tracker/internal/telegramclient"

type T struct {
	telegramclient.Config `validate:"required"`
}
