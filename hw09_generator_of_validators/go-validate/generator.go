package main

import (
	"errors"
	"github.com/balabanovds/otus-golang/hw09_generator_of_validators/pkg/parser"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrFieldLen   = errors.New("validate: field does not match length")
	ErrFieldRegex = errors.New("validate: field does not match regex")
	ErrFieldIn    = errors.New("validate: field does not exists in variety")
	ErrFieldMin   = errors.New("validate: field is less than min value")
	ErrFieldMax   = errors.New("validate: field is greater than max value")
)

type ValidationError struct {
	fieldName string
	err       error
}

type User struct {
	ID     string `json:"id" validate:"len:36"`
	Name   string
	Age    int      `validate:"min:18|max:50"`
	Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	Role   string   `validate:"in:admin,stuff"`
	Num    int      `validate:"in:1,2"`
	Phones []string `validate:"len:11"`
}

func (u User) Validate() ([]ValidationError, error) {
	var vErrors []ValidationError

	// processing 'len' tag
	i, err := strconv.Atoi("36")
	if err != nil {
		return nil, err
	}

	if len(u.ID) != i {
		vErrors = append(vErrors, ValidationError{
			fieldName: "ID",
			err:       ErrFieldLen,
		})
	}

	// processing 'min' tag
	i, err = strconv.Atoi("18")
	if err != nil {
		return nil, err
	}
	if u.Age < i {
		vErrors = append(vErrors, ValidationError{
			fieldName: "Age",
			err:       ErrFieldMin,
		})
	}

	// processing 'regex' tag
	r := regexp.MustCompile("^\\w+@\\w+\\.\\w+$")
	if !r.MatchString(u.Email) {
		vErrors = append(vErrors, ValidationError{
			fieldName: "Email",
			err:       ErrFieldRegex,
		})
	}

	// processing 'in' tag
	// for strings
	var in bool
	for _, s := range strings.Split("admin,stuff", ",") {
		if u.Role == s {
			in = true
			break
		}
	}
	if !in {
		vErrors = append(vErrors, ValidationError{
			fieldName: "Email",
			err:       ErrFieldIn,
		})
	}

	return vErrors, nil
}

func generate(filepath string, data *parser.ParsedData) error {

	return nil
}
