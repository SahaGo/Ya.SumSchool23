package controllers

import (
	"github.com/go-playground/validator/v10"
	"regexp"
	"strings"
	"time"
)

type CustomValidator struct {
	validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()
	_ = v.RegisterValidation("hh_mm_interval", IsHhMmInterval)

	return &CustomValidator{
		validator: v,
	}
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		return err
	}
	return nil
}

func IsHhMmInterval(fl validator.FieldLevel) bool {
	str := fl.Field().String()

	r := regexp.MustCompile("^\\d\\d:\\d\\d-\\d\\d:\\d\\d$")
	if !r.MatchString(str) {
		return false
	}

	layout := "15:04"
	left, right, _ := strings.Cut(str, "-")
	leftTime, err := time.Parse(layout, left)
	if err != nil {
		return false
	}

	rightTime, err := time.Parse(layout, right)
	if err != nil {
		return false
	}

	if !rightTime.After(leftTime) {
		return false
	}

	return true
}
