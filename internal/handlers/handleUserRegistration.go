package handlers

import (
	"chrispaul1/chirpy/internal/database"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func (h *UserHandler) HandleUserRegistration(w http.ResponseWriter, req *http.Request) {
	type User struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	type userEmail struct {
		Body string `json:"email"`
	}
	decoder := json.NewDecoder(req.Body)
	newEmail := userEmail{}
	err := decoder.Decode(&newEmail)
	if err != nil {
		errMsg := "Something went wrong"
		RespondWithError(w, 400, errMsg)
	}

	userStruct := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Email:     newEmail.Body,
	}

	newUser, err := h.cfg.DB.CreateUser(req.Context(), userStruct)
	bodyUser := User{
		ID:        newUser.ID,
		CreatedAt: newUser.CreatedAt,
		UpdatedAt: newUser.UpdatedAt,
		Email:     newUser.Email,
	}
	if err != nil {
		errMsg := "Something went wrong in creating the user"
		RespondWithError(w, 400, errMsg)
	}
	RespondWithJson(w, 201, bodyUser)
}
