package utils

import (
	"strconv"

	"github.com/google/uuid"
)

func GetInt(val string) *int {
	if v, err := strconv.Atoi(val); err == nil {
		return &v
	}
	return nil
}

// IsValidUUID checks if a string is a valid UUID
func IsValidUUID(id string) bool {
	if id == "" {
		return false
	}

	_, err := uuid.Parse(id)
	return err == nil
}
