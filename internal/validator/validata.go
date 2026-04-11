package validator

import (
	"net/url"
	"regexp"
)

var KeyRX =  regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func (v *Validator) IsURL (path string) bool {
	u, err := url.ParseRequestURI(path)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func (v *Validator) CheckLength(value string, min, max int, key, message string) {
    v.Check(len(value) >= min && len(value) <= max, key, message)
}