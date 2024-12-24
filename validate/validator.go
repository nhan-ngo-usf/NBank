package validate

import (
	"fmt"
	"net/mail"
	"regexp"
	"strconv"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullName = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)
func ValidateString(value string, minLength int, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("must contain from %d-%d characters", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("must contain only lowercase letters, digits or underscore")
	}
	return nil
}

func ValidateFullName(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if !isValidFullName(value) {
		return fmt.Errorf("must contain only letters or spaces")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 20)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 100); err != nil {
		return err
	}
	if _, err := mail.ParseAddress((value)); err != nil {
		return fmt.Errorf("is not a valid email address")
	}
	return nil
}

func ValidateBalance(value int64) error {
	if value < 0 {
		return fmt.Errorf("account balance cannot be negative")
	}
	return nil
}

func contain(list []string, v string) bool {
	for _, value := range list {
		if value == v {
			return true
		}
	}
	return false
}

func ValidateCurrency(value string) error {
	var CurrencyList = []string {"USD", "EUR", "POUND"}
	if !contain(CurrencyList, value) {
		return fmt.Errorf("unsupported currency, list of supported currency %v", CurrencyList)
	}
	return nil
}

func ValidateAccountID(value string) error {
	id, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("cannot convert id into integer %s", err)
	}
	if id < 0 {
		return fmt.Errorf("invalid id value %s", err)
	}
	return nil
}

func ValidateVerifyCode(value string) error {
	err := ValidateString(value, 32,128)
	fmt.Printf("invalid value %s", value)

	return err
}

func ValidateEmailID(value int64) error {
	if value <= 0 {
		return fmt.Errorf("must be a positive integer")
	}
	return nil
}