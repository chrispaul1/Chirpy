package handlers

import (
	"chrispaul1/chirpy/internal/auth"
	"chrispaul1/chirpy/internal/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *UserHandler) HandleUserRegistration(w http.ResponseWriter, req *http.Request) {
	type User struct {
		ID          uuid.UUID `json:"id"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
		Email       string    `json:"email"`
		IsChirpyRed bool      `json:"is_chirpy_red"`
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

	userHashedPass, err := auth.HashPassword(userFields.Password)
	if err != nil {
		errMsg := "Password not accepted"
		RespondWithError(w, 400, errMsg)
		return
	}

	userStruct := database.CreateUserParams{
		ID:             uuid.New(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
		Email:          userFields.Email,
		HashedPassword: userHashedPass,
		IsChirpyRed:    false,
	}

	newUser, err := h.cfg.DB.CreateUser(req.Context(), userStruct)
	if err != nil {
		//errMsg := "Something went wrong in creating the user"
		RespondWithError(w, 400, err.Error())
		return
	}

	userResponse := User{
		ID:          newUser.ID,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed,
	}
	RespondWithJson(w, 201, userResponse)
}
