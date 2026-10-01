package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

	type isValid struct {
		Valid bool `json:"valid"`
	}

	var str toValidate
	decoder := json.NewDecoder(req.Body)
	err := decoder.Decode(&str)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if len(str.Body) <= 140 {
		respondWithJSON(w, 200, isValid{Valid: true})
	} else {
		respondWithError(w, 400, "Chirp is too long")
	}
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
