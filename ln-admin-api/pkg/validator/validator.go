package validator

import (
	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Validate 验证结构体
func Validate(s interface{}) error {
	return validate.Struct(s)
}

// RegisterValidation 注册自定义验证器
func RegisterValidation(tag string, fn validator.Func) error {
	return validate.RegisterValidation(tag, fn)
}
