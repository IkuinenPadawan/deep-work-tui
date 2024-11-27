package utils

import "time"

func ParseTime(t string) time.Time {
	parsed, _ := time.Parse("15:04", t)
	return parsed
}

func IsValidTime(t string) bool {
	_, err := time.Parse("15:04", t)
	return err == nil
}
