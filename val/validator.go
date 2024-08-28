package val

import (
	"fmt"
	"net/mail"
	"regexp"

	"github.com/CodeSingerGnC/MicroBank/otpcode"
)

var (
	isValidUserAccount = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidUsername    = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
	isValidPasscode    = regexp.MustCompile(`^[0-9]+$`).MatchString
)

func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUserAccount(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUserAccount(value) {
		return fmt.Errorf("must contain only lowercase letters, digits, or underscore")
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 200); err != nil {
		return err
	}
	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	return nil
}

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}

func ValidatePasscode(value string) error {
	if len(value) != int(otpcode.Digits) {
		return fmt.Errorf("passcode must be %d digits", otpcode.Digits)
	}
	if !isValidUserAccount(value) {
		return fmt.Errorf("passcode must only contain digit")
	}
	return nil
}
