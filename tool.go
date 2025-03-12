package main

import (
	"regexp"
	"strings"
)

func convertCamelToSnake(s string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`) // Match lowercase + uppercase transition
	snake := re.ReplaceAllString(s, `${1}_${2}`)  // Insert underscore between matches
	return strings.ToLower(snake)                 // Convert to lowercase
}
