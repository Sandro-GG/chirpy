package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Sandro-GG/chirpy/internal/auth"
	"github.com/Sandro-GG/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	params := &parameters{}

	err := json.NewDecoder(req.Body).Decode(params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	pwd, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: pwd,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Error creating user", err)
		return
	}

	responseUser := User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusCreated, responseUser)
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, req *http.Request) {
	type parameters struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	params := &parameters{}

	err := json.NewDecoder(req.Body).Decode(params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	dbUser, err := cfg.db.GetUserByEmail(req.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	if ok, err := auth.CheckPasswordHash(params.Password, dbUser.HashedPassword); !ok || err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	duration := 1 * time.Hour
	now := time.Now().UTC()

	token, err := auth.MakeJWT(dbUser.ID, cfg.secret, duration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to make JWT", err)
		return
	}

	refToken := auth.MakeRefreshToken()
	_, err = cfg.db.CreateRefreshToken(req.Context(), database.CreateRefreshTokenParams{
		Token:     refToken,
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    dbUser.ID,
		ExpiresAt: now.Add(60 * 24 * time.Hour),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't save refresh token", err)
		return
	}

	type response struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}

	respondWithJSON(w, http.StatusOK, response{
		ID:           dbUser.ID,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		Email:        dbUser.Email,
		Token:        token,
		RefreshToken: refToken,
		IsChirpyRed:  dbUser.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, req *http.Request) {
	refToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	tokenRow, err := cfg.db.GetUserFromRefreshToken(req.Context(), refToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	now := time.Now().UTC()
	if tokenRow.ExpiresAt.Before(now) {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", errors.New("refresh token expired"))
		return
	}

	if tokenRow.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", errors.New("refresh token revoked"))
		return
	}

	tokenStr, err := auth.MakeJWT(tokenRow.ID, cfg.secret, 1*time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to make JWT", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		Token: tokenStr,
	})
}

func (cfg *apiConfig) handlerRevokeRefreshToken(w http.ResponseWriter, req *http.Request) {
	refToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	now := time.Now().UTC()

	err = cfg.db.RevokeRefreshToken(req.Context(), database.RevokeRefreshTokenParams{
		Token: refToken,
		RevokedAt: sql.NullTime{
			Time:  now,
			Valid: true,
		},
		UpdatedAt: now,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke refresh token", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handlerUpdateEmailPassword(w http.ResponseWriter, req *http.Request) {
	token, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	params := &parameters{}

	err = json.NewDecoder(req.Body).Decode(params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	userUUID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	hash, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to hash password", err)
		return
	}

	updatedUser, err := cfg.db.UpdateUser(req.Context(), database.UpdateUserParams{
		ID:             userUUID,
		Email:          params.Email,
		HashedPassword: hash,
		UpdatedAt:      time.Now().UTC(),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to update user", err)
		return
	}

	user := &User{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, user)
}

func (cfg *apiConfig) handlerUpgradeUser(w http.ResponseWriter, req *http.Request) {
	apiKey, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unable to check API key", err)
		return
	}

	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "API keys don't match", nil)
		return
	}

	type parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	params := &parameters{}
	err = json.NewDecoder(req.Body).Decode(params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Unable to decode JSON", err)
		return
	}

	if params.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpgradeUser(req.Context(), params.Data.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		respondWithError(w, http.StatusNotFound, "User was not found", err)
		return
	}
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Unable to upgrade user", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
