package request

import (
	"github.com/go-playground/mold/v4"
	"github.com/go-playground/mold/v4/modifiers"
	validateV10 "github.com/go-playground/validator/v10"
)

var (
	validator   *validateV10.Validate
	transformer *mold.Transformer
)

func init() {
	validator = validateV10.New()
	transformer = modifiers.New()
}
