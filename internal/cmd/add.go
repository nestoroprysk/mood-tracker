package cmd

import "fmt"

func newAdd(args ...string) (Cmd, error) {
	return func() (string, error) {
		return fmt.Sprintf("%v", args), nil
	}, nil
}
