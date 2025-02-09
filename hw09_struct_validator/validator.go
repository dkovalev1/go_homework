package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// ValidationError represents a single validation error.
type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Err)
}

func (v ValidationError) Unwrap() error { return v.Err }

type SyntaxError struct {
	Field string
	Err   error
}

func (s SyntaxError) Error() string {
	return fmt.Sprintf("%s: Syntax error: %s", s.Field, s.Err)
}

func (s SyntaxError) Unwrap() error { return s.Err }

func NewSyntaxError(field string, err error) *SyntaxError {
	return &SyntaxError{
		Field: field,
		Err:   err,
	}
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	sb := strings.Builder{}

	for i, err := range v {
		if i != 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("%s: %s", err.Field, err.Err))
	}
	return sb.String()
}

func (v ValidationErrors) Unwrap() error {
	return v[0]
}

func NewValidationErrors() ValidationErrors {
	return make(ValidationErrors, 0)
}

type IRule interface {
	Validate(fld reflect.StructField, v reflect.Value) *ValidationError
}

type Rule struct {
	Command string
	Args    []string
}

type StringRule struct {
	Rule
}

type IntRule struct {
	Rule
}

var (
	ErrInvalidRule = errors.New("invalid rule")
	ErrValidation  = errors.New("validation error")
)

func NewStringRule(rule string) (IRule, error) {
	parts := strings.SplitN(rule, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidRule, rule)
	}
	ret := StringRule{
		Rule: Rule{
			Command: parts[0],
			Args:    strings.Split(parts[1], ","),
		},
	}
	return ret, nil
}

func NewIntRule(rule string) (IRule, error) {
	parts := strings.SplitN(rule, ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidRule, rule)
	}
	ret := IntRule{
		Rule: Rule{
			Command: parts[0],
			Args:    strings.Split(parts[1], ","),
		},
	}
	return ret, nil
}

func makeRule(fieldType reflect.Type, ruleText string) (IRule, error) {
	kind := fieldType.Kind()

	var rule IRule
	var err error

	switch kind { //nolint:exhaustive
	case reflect.String:
		rule, err = NewStringRule(ruleText)

	case reflect.Int:
		rule, err = NewIntRule(ruleText)

	case reflect.Array:
		aType := reflect.ArrayOf(1, fieldType)
		return makeRule(aType, ruleText)

	default:
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidRule, ruleText)
	}

	return rule, nil
}

func (s StringRule) validateStringValue(stringVal string) (err error) {
	switch s.Command {
	case "len":
		requiredLen, err := strconv.Atoi(s.Args[0])
		if err != nil {
			return err
		}
		if len(stringVal) != requiredLen {
			return fmt.Errorf(
				"%w value '%s': required len: %d, actual len: %d",
				ErrValidation,
				stringVal, requiredLen, len(stringVal))
		}
		return nil

	case "regexp":
		result, err := regexp.Match(s.Args[0], []byte(stringVal))
		if err != nil {
			return err
		}
		if !result {
			return fmt.Errorf("%w value '%s': regexp not matched",
				ErrValidation, stringVal)
		}
		return nil

	case "in":
		for _, arg := range s.Args {
			if arg == stringVal {
				return nil
			}
		}
		return fmt.Errorf("%w value '%s': not in: %s",
			ErrValidation, stringVal, strings.Join(s.Args, ","))
	}

	return nil
}

func (s StringRule) Validate(field reflect.StructField, v reflect.Value) *ValidationError {
	var err error

	if field.Type.Kind() == reflect.Array {
		for i := 0; i < v.Len(); i++ {
			elem := v.Index(i).String()

			err = s.validateStringValue(elem)
			if err != nil {
				break
			}
		}
	} else {
		err = s.validateStringValue(v.String())
	}

	if err != nil {
		return &ValidationError{
			Field: field.Name,
			Err:   err,
		}
	}

	return nil
}

func (i IntRule) validateIntValue(intVal int64) error {
	switch i.Command {
	case "min":
		requiredMin, err := strconv.Atoi(i.Args[0])
		if err != nil {
			return err
		}
		if intVal < int64(requiredMin) {
			return fmt.Errorf("min: %d", requiredMin)
		}

	case "max":
		requiredMax, err := strconv.Atoi(i.Args[0])
		if err != nil {
			return err
		}
		if intVal > int64(requiredMax) {
			return fmt.Errorf("max: %d", requiredMax)
		}
	case "in":
		for _, arg := range i.Args {
			argInt, err := strconv.Atoi(arg)
			if err != nil {
				return err
			}
			if intVal == int64(argInt) {
				return nil
			}
		}
		return fmt.Errorf("%w value '%d': not in: %s",
			ErrValidation,
			intVal, strings.Join(i.Args, ","))
	}
	return nil
}

// Validate implements IRule.
func (i IntRule) Validate(field reflect.StructField, v reflect.Value) *ValidationError {
	var err error
	if field.Type.Kind() == reflect.Array {
		for idx := 0; idx < v.Len(); idx++ {
			elem := v.Index(idx).Int()

			err = i.validateIntValue(elem)
			if err != nil {
				break
			}
		}
	} else {
		err = i.validateIntValue(v.Int())
	}

	if err != nil {
		return &ValidationError{
			Field: field.Name,
			Err:   err,
		}
	}

	return nil
}

func Validate(v interface{}) error {
	// Place your code here.
	st := reflect.TypeOf(v)
	ret := NewValidationErrors()

	if st.Kind() != reflect.Struct {
		ret = append(ret, ValidationError{
			Field: "<input>",
			Err:   errors.New("not a struct"),
		})
		return ret
	}

	for fieldNo := 0; fieldNo < st.NumField(); fieldNo++ {
		field := st.Field(fieldNo)
		if !field.IsExported() {
			continue
		}
		tag := field.Tag.Get("validate")
		value := reflect.ValueOf(v).Field(fieldNo)

		if tag != "" {
			rules := strings.Split(tag, "|")

			ruleSet := make([]IRule, 0, len(rules))

			// Check the rules to make sure all rules for this field are valid
			for _, rule := range rules {
				irule, syntaxErr := makeRule(field.Type, rule)
				if syntaxErr != nil {
					return NewSyntaxError(field.Name, syntaxErr)
				}

				if irule == nil {
					continue
				}
				ruleSet = append(ruleSet, irule)
			}

			// Apply the rules when we know all rules are valid
			for _, irule := range ruleSet {
				verr := irule.Validate(field, value)
				if verr != nil {
					ret = append(ret, *verr)
				}
			}
		}
	}
	if len(ret) > 0 {
		return ret
	}
	return nil
}
