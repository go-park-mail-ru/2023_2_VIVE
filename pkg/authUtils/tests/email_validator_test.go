package authUtils

import (
	"HnH/pkg/authUtils"
	"HnH/pkg/testHelper"
	"fmt"
	"testing"

	emailverifier "github.com/AfterShip/email-verifier"
)

var testValidateEmailCases = []struct {
	email       string
	expectedErr error
}{
	{
		email:       "email@example.com",
		expectedErr: nil,
	},
	{
		email:       "email@mail.ru",
		expectedErr: nil,
	},
	{
		email:       "email@yandex.ru",
		expectedErr: nil,
	},
	{
		email:       "email@google.ru",
		expectedErr: nil,
	},
	{
		email:       "email@google.com",
		expectedErr: nil,
	},
	// {
	// 	email:       "email@google.co",
	// 	expectedErr: authUtils.INVALID_EMAIL,
	// },
	{
		email:       "email@mail",
		expectedErr: authUtils.INVALID_EMAIL,
	},
	{
		email:       "email mail.ru",
		expectedErr: authUtils.INVALID_EMAIL,
	},
	{
		email:       "mail.ru",
		expectedErr: authUtils.INVALID_EMAIL,
	},
	{
		email:       "mail.ru",
		expectedErr: authUtils.INVALID_EMAIL,
	},
	{
		email:       "",
		expectedErr: authUtils.EMPTY_EMAIL,
	},
}

func TestValidateEmail(t *testing.T) {
	for _, testCase := range testValidateEmailCases {
		actualErr := authUtils.ValidateEmail(testCase.email)
		if actualErr != testCase.expectedErr && actualErr != fmt.Errorf(emailverifier.ErrTimeout) {
			t.Errorf(testHelper.ErrNotEqual(testCase.expectedErr, actualErr))
		}
	}
}
