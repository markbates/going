package validate

import (
	"sync"

	"github.com/markbates/going/wait"
)

// ValidationErrors holds onto all of the error messages
// that get generated during the validation process.
type ValidationErrors struct {
	Errors map[string][]string `json:"errors"`
	Lock   *sync.RWMutex       `json:"-"`
}

// Validator must be implemented in order to pass the
// validator object into the Validate function.
type Validator interface {
	IsValid(errors *ValidationErrors)
}

// NewValidationErrors returns a pointer to a ValidationErrors
// object that has been primed and ready to go.
func NewValidationErrors() *ValidationErrors {
	return &ValidationErrors{
		Errors: make(map[string][]string),
		Lock:   new(sync.RWMutex),
	}
}

// Count returns the number of errors.
func (v *ValidationErrors) Count() int {
	return len(v.Errors)
}

// HasAny returns true/false depending on whether any errors
// have been tracked.
func (v *ValidationErrors) HasAny() bool {
	return v.Count() > 0
}

// Append concatenates two ValidationErrors objects together.
// This will modify the first object in place.
func (v *ValidationErrors) Append(ers *ValidationErrors) {
	for key, value := range ers.Errors {
		for _, msg := range value {
			v.Add(key, msg)
		}
	}
}

// Add will add a new message to the list of errors using
// the given key. If the key already exists the message will
// be appended to the array of the existing messages.
func (v *ValidationErrors) Add(key string, msg string) {
	v.Lock.Lock()
	v.Errors[key] = append(v.Errors[key], msg)
	v.Lock.Unlock()
}

// Get returns an array of error messages for the given key.
func (v *ValidationErrors) Get(key string) []string {
	return v.Errors[key]
}

// Validate takes in n number of Validator objects and will run
// them and return back a point to a ValidationErrors object that
// will contain any errors.
func Validate(validators ...Validator) *ValidationErrors {
	errors := NewValidationErrors()

	wait.Wait(len(validators), func(index int) {
		validator := validators[index]
		validator.IsValid(errors)
	})

	return errors
}
