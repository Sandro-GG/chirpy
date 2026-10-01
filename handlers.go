package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := fmt.Sprintf(`<html>
  		<body>
    		<h1>Welcome, Chirpy Admin</h1>
    		<p>Chirpy has been visited %d times!</p>
  		</body>
	</html>`, cfg.fileserverHits.Load())
	w.Write([]byte(html))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, req *http.Request) {
	type toValidate struct {
		Body string `json:"body"`
	}

	type cleaned struct {
		CleanedBody string `json:"cleaned_body"`
	}

	var str toValidate
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&str)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	badWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	if len(str.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}

	cleanBody := getCleanedBody(str.Body, badWords)

	respondWithJSON(w, 200, cleaned{CleanedBody: cleanBody})
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type toError struct {
		Error string `json:"error"`
	}

	errToReturn := toError{
		Error: msg,
	}

	dat, err := json.Marshal(errToReturn)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	res, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(res)
}

func getCleanedBody(text string, badWords map[string]struct{}) string {
	words := strings.Fields(text)
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if _, ok := badWords[loweredWord]; ok {
			words[i] = "****"
		}
	}

	return strings.Join(words, " ")
}
