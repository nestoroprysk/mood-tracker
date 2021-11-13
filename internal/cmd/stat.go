package cmd

import "fmt"

func newStat(args ...string) (Cmd, error) {
	return func() (string, error) {
		return fmt.Sprintf("%v", args), nil
	}, nil
}
