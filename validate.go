package main

import "strings"

var badWords = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func getCleanedBody(text string, badWords map[string]struct{}) string {
	words := strings.Split(text, " ")
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
