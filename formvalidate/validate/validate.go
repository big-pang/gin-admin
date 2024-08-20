package validate

import (
	"github.com/go-playground/validator/v10"
	"sync"
)

var (
	validate *validator.Validate
	once     sync.Once
)

func Init() {
	once.Do(func() {
		validate = validator.New()
	})
}

func Struct(s interface{}) error {
	Init()
	return validate.Struct(s)
}
