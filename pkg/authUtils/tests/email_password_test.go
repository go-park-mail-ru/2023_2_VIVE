package authUtils

import (
	"HnH/pkg/authUtils"
	"HnH/pkg/testHelper"
	"strings"
	"testing"
)

func tooLongPassword() string {
	collectedSymbols := []string{"A", "b", "1"}
	symbols := make([]string, 129)

	for i := range symbols {
		symbols[i] = collectedSymbols[i%len(collectedSymbols)]
	}

	return strings.Join(symbols, "")
}

var testValidatePasswordCases = []struct {
	password    string
	expectedErr error
}{
	{
		password:    "Qwerty123",
		expectedErr: nil,
	},
	{
		password:    "ASDadfasdf1435",
		expectedErr: nil,
	},
	{
		password:    "Good_Password_55",
		expectedErr: nil,
	},
	{
		password:    "Qq345",
		expectedErr: authUtils.INVALID_PASSWORD,
	},
	{
		password:    tooLongPassword(),
		expectedErr: authUtils.INVALID_PASSWORD,
	},
	{
		password:    "NoDigits",
		expectedErr: authUtils.INVALID_PASSWORD,
	},
	{
		password:    "no_capital123",
		expectedErr: authUtils.INVALID_PASSWORD,
	},
	{
		password:    "",
		expectedErr: authUtils.EMPTY_PASSWORD,
	},
}

func TestValidatePassword(t *testing.T) {
	for _, testCase := range testValidatePasswordCases {
		actualErr := authUtils.ValidatePassword(testCase.password)
		if actualErr != testCase.expectedErr {
			t.Errorf(testHelper.ErrNotEqual(testCase.expectedErr, actualErr))
		}
	}
}
