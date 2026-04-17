package utils

import (
	"errors"

	"github.com/google/uuid"
)

var ErrGenerateUUID = errors.New("failed to generate uuid")

func GenerateUUID() (string, error) {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		return "", ErrGenerateUUID
	}

	return uuidV7.String(), nil
}
