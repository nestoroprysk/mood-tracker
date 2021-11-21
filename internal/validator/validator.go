package validator

import "gopkg.in/go-playground/validator.v8"

// Validator validates any struct.
type Validator interface {
	// Struct validates a structure.
	Struct(s interface{}) error
}

var _ Validator = &validate{}

type validate struct {
	*validator.Validate
}

// New creates a validator that is defined by the `validate` tag.
func New() Validator {
	conf := &validator.Config{TagName: "validate"}
	return &validate{Validate: validator.New(conf)}
}
