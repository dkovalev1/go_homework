package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require" //nolint:all
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

var err1 = ValidationErrors{
	ValidationError{
		Field: "ID",
		Err:   errors.New("validation error value '123e4567-e89b-12d3-a456-4266554400000': required len: 36, actual len: 37"),
	},
	ValidationError{
		Field: "Email",
		Err:   errors.New("validation error value 'RQ3nX@example.co!m': regexp not matched"),
	},
	ValidationError{
		Field: "Role",
		Err:   errors.New("validation error value 'visitor': not in: admin,stuff"),
	},
}

func TestValidateNoStruct(t *testing.T) {
	err := Validate(1)
	require.Error(t, err)
	require.EqualError(t, err, ValidationErrors{
		ValidationError{
			Field: "<input>",
			Err:   errors.New("not a struct"),
		},
	}.Error())

	require.EqualError(t, Validate(1), "<input>: not a struct")
}

func TestValidateError(t *testing.T) {
	err := Validate(User{})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrValidation)

	// Check that composed list of errors is equal to expected string
	require.EqualError(t, err, `ID: validation error value '': required len: 36, actual len: 0
Age: min: 18
Email: validation error value '': regexp not matched
Role: validation error value '': not in: admin,stuff`)
}

func TestValidateSuccess(t *testing.T) {
	err := Validate(User{
		ID:    "123e4567-e89b-12d3-a456-426655440000",
		Name:  "Bob",
		Age:   20,
		Email: "RQ3nX@example.com",
		Role:  "admin",
	})
	require.NoError(t, err)
}

func TestValidateSyntaxError(t *testing.T) {
	type Admin struct {
		User
		Login string `validate:"len-5"`
	}

	err := Validate(Admin{})
	require.Error(t, err)
	require.ErrorIs(t, err, ErrInvalidRule)
	require.EqualError(t, err, "Login: Syntax error: invalid rule: len-5")
}

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:    "123e4567-e89b-12d3-a456-426655440000",
				Name:  "Bob",
				Age:   20,
				Email: "RQ3nX@example.com",
				Role:  "admin",
			},
			expectedErr: nil,
		},

		{
			/* invalid ID, role and email */
			in: User{
				ID:    "123e4567-e89b-12d3-a456-4266554400000",
				Name:  "Bob",
				Age:   20,
				Email: "RQ3nX@example.co!m",
				Role:  "visitor",
				Phones: []string{
					"+380501234567",
					"+380503456789",
				},
			},
			expectedErr: err1,
		},

		{
			in: App{
				Version: "1.0.0",
			},
			expectedErr: nil,
		},

		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},

		{
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},

		{
			in: Response{
				Code: 201,
				Body: "OK",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   errors.New("validation error value '201': not in: 200,404,500"),
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)

			if tt.expectedErr == nil {
				require.NoError(t, err)
				return
			}

			require.Error(t, err)
			require.ErrorIs(t, err, ErrValidation)

			var expected, validated ValidationErrors

			errors.As(tt.expectedErr, &expected)
			errors.As(err, &validated)

			require.Equal(t, len(expected), len(validated), "Validate failed")
			for i, v := range validated {
				require.ErrorIs(t, v.Err, ErrValidation)

				require.Equal(t, expected[i].Field, v.Field)
				require.Equal(t, expected[i].Err.Error(), v.Err.Error())
			}
		})
	}
}
