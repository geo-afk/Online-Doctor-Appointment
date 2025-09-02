package mailing

import (
	"errors"

	"github.com/AfterShip/email-verifier"
)

var (
	verifier = emailverifier.NewVerifier()

	ErrInvalidEmail = errors.New("invalid email address")
)

func EmailChecker(email string) (bool, error) {

	ret, err := verifier.Verify(email)
	if err != nil {
		return false, err
	}

	if !ret.Syntax.Valid || !ret.HasMxRecords || ret.Disposable || ret.RoleAccount {
		return false, ErrInvalidEmail
	} else {

		return true, nil
	}
}
