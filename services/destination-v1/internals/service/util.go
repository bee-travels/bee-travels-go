package service

import "strings"

func normalize(location string) string {
	split := strings.Split(strings.ToLower(location), " ")
	return strings.Join(split, "-")
}
