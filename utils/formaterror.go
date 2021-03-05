package utils

import (
	"errors"
	"strings"
)

func FormatError(err string) error {
	if strings.Contains(err, "nickname") && strings.Contains(err, "email") {
		return errors.New("Nickname or Email already Taken")
	}

	if strings.Contains(err, "nickname") {
		return errors.New("Nickname Already Taken")
	}

	if strings.Contains(err, "email") {
		return errors.New("Email already taken")
	}

	if strings.Contains(err, "title") {
		return errors.New("Title Already Taken")
	}

	if strings.Contains(err, "hashedPassword") {
		return errors.New("Incorrect password")
	}

	return errors.New("Incorrect details")
}
