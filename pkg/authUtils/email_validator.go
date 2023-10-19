package authUtils

import (
	emailverifier "github.com/AfterShip/email-verifier"
)

var verifier = emailverifier.NewVerifier()

func ValidateEmail(email string) error {
	if len(email) == 0 {
		return EMPTY_EMAIL
	}

	ret, err := verifier.Verify(email)
	if err != nil {
		return err
	} else if !ret.Syntax.Valid {
		return INVALID_EMAIL
	}

	return nil
}
