package handlers

import (
	"chrispaul1/chirpy/internal/auth"
	"chrispaul1/chirpy/internal/database"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *UserHandler) HandleUserLogin(w http.ResponseWriter, req *http.Request) {

	type userResponseStruct struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
		IsChirpyRed  bool      `json:"is_chirpy_red"`
	}

	type userParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	userFields := userParams{}
	err := decoder.Decode(&userFields)

	if err != nil {
		errMsg := "Something went wrong"
		RespondWithError(w, 400, errMsg)
		return
	}

	user, err := h.cfg.DB.GetUserByEmail(req.Context(), userFields.Email)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	err = auth.CheckPasswordHash(userFields.Password, user.HashedPassword)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	jwtToken, err := auth.MakeJWT(user.ID, h.cfg.JWT_SECRET)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, "Error, JWT token could not created")
		return
	}

	newRefreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	refreshTokenStruct := database.CreateRefreshTokenParams{
		Token:     newRefreshToken,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
		RevokedAt: sql.NullTime{
			Valid: false,
		},
	}

	_, err = h.cfg.DB.CreateRefreshToken(req.Context(), refreshTokenStruct)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userResponse := userResponseStruct{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        jwtToken,
		RefreshToken: newRefreshToken,
		IsChirpyRed:  user.IsChirpyRed,
	}
	RespondWithJson(w, http.StatusOK, userResponse)
}
