package main

import "strings"

func capitalize(name string) string {
	name = strings.ToLower(name)
	tokens := strings.Split(name, "-")
	for i, token := range tokens {
		if includes(token) {
			tokens[i] = strings.ToUpper(string(token[0])) + token[1:]
		}
	}
	joiner := " "
	if strings.Contains(name, lowercaseExceptions[2]) {
		joiner = "-"
	}
	return strings.Join(tokens, joiner)
}

func includes(token string) bool {
	for _, exception := range lowercaseExceptions {
		if exception == token {
			return true
		}
	}
	return false
}
