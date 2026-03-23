package utils

import "strconv"

func GetInt(val string) *int {
	if v, err := strconv.Atoi(val); err == nil {
		return &v
	}
	return nil
}
