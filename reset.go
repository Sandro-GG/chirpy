package main

import "net/http"

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, req *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Forbidden", nil)
		return
	}

	err := cfg.db.DeleteUsers(req.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error deleting users", err)
		return
	}

	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}
