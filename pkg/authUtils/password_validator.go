package authUtils

import "unicode"

const specialChars = `~!?@#$%^&*_-+()[]{}></\|"'.,:;`

func ValidatePassword(password string) error {
	if len(password) == 0 {
		return EMPTY_PASSWORD
	}

	if len(password) < 8 || len(password) > 128 {
		return INVALID_PASSWORD
	}

	hasDigit := false
	hasCapital := false
	hasSpecialChar := false

	for _, char := range password {
		if !hasDigit {
			hasDigit = unicode.IsDigit(char)
		}

		if !hasCapital {
			hasCapital = unicode.Is(unicode.Latin, char) && unicode.IsUpper(char)
		}

		if !hasSpecialChar {
			for _, specCh := range specialChars {
				hasSpecialChar = (char == specCh)

				if hasSpecialChar {
					break
				}
			}
		}
	}

	if !(hasDigit && hasCapital && hasSpecialChar) {
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
