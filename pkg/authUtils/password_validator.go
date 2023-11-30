package authUtils

import "unicode"

// const specialChars = `~!?@#$%^&*_-+()[]{}></\|"'.,:;`

func ValidatePassword(password string) error {
	if len(password) == 0 {
		return EMPTY_PASSWORD
	}

	if len(password) < 6 || len(password) > 128 {
		return INVALID_PASSWORD
	}

	hasDigit := false
	hasCapital := false
	// hasSpecialChar := false

	for _, char := range password {
		if !hasDigit {
			hasDigit = unicode.IsDigit(char)
		}

		if !hasCapital {
			hasCapital = unicode.Is(unicode.Latin, char) && unicode.IsUpper(char)
		}
	}

	if !(hasDigit && hasCapital /* && hasSpecialChar */) {
		return INVALID_PASSWORD
	}

	return nil
}

func IsPasswordEmpty(password string) error {
	if len(password) == 0 {
		return EMPTY_PASSWORD
	}

	return nil
}
