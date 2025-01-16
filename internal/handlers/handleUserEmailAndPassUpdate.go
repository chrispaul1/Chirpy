package handlers

import (
	"chrispaul1/chirpy/internal/auth"
	"chrispaul1/chirpy/internal/database"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (h *UserHandler) HandleUserEmailAndPassUpdate(w http.ResponseWriter, req *http.Request) {
	//user struct
	type User struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
	}

	//Decode the user password and email from request
	type userRequest struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	decoder := json.NewDecoder(req.Body)
	userFields := userRequest{}
	err := decoder.Decode(&userFields)

	if err != nil {
		errMsg := "Something went wrong"
		RespondWithError(w, http.StatusUnauthorized, errMsg)
		return
	}

	//retrieve the jwt token from the header
	authHeader := req.Header.Get("Authorization")
	if authHeader == "" {
		RespondWithError(w, http.StatusUnauthorized, "Error,Header not found")
		return
	}

	jwtToken, ok := strings.CutPrefix(authHeader, "Bearer ")
	if !ok {
		RespondWithError(w, http.StatusUnauthorized, "Token could not be found")
		return
	}

	userID, err := auth.ValidateJWT(jwtToken, h.cfg.JWT_SECRET)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	newHashedPassword, err := auth.HashPassword(userFields.Password)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userEmailPassStruct := database.UpdateUserEmailAndPasswordParams{
		Email:          userFields.Email,
		HashedPassword: newHashedPassword,
		UpdatedAt:      time.Now(),
		ID:             userID,
	}

	updatedUser, err := h.cfg.DB.UpdateUserEmailAndPassword(req.Context(), userEmailPassStruct)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	userResponse := User{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAt:   updatedUser.UpdatedAt,
		Email:       updatedUser.Email,
		IsChirpyRed: updatedUser.IsChirpyRed,
	}

	RespondWithJson(w, http.StatusOK, userResponse)
}
