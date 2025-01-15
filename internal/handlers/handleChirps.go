package handlers

import (
	"chrispaul1/chirpy/internal/auth"
	"chrispaul1/chirpy/internal/database"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (h *UserHandler) HandleChirps(w http.ResponseWriter, req *http.Request) {
	type chirpBody struct {
		Body string `json:"body"`
	}

	tokenString, err := auth.GetBearerToken(req.Header)

	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	jwtUserID, err := auth.ValidateJWT(tokenString, h.cfg.JWT_SECRET)
	if err != nil {
		RespondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	decoder := json.NewDecoder(req.Body)
	userChirp := chirpBody{}
	err = decoder.Decode(&userChirp)

	if err != nil {
		errMsg := "Something went wrong"
		RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	if len(userChirp.Body) > 140 {
		errMsg := "Chirp is too long"
		RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	if len(userChirp.Body) == 0 || userChirp.Body == "" {
		errMsg := "Error, chirp cannot be empty"
		RespondWithError(w, http.StatusBadRequest, errMsg)
		return
	}

	cleanedBody := cleanText(userChirp.Body)

	chirpStruct := database.CreateChirpParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Body:      cleanedBody,
		UserID:    jwtUserID,
	}
	newChirp, err := h.cfg.DB.CreateChirp(req.Context(), chirpStruct)
	if err != nil {
		errMsg := "Error, could not insert chirp into database"
		RespondWithError(w, 400, errMsg)
		return
	}

	jsonChrip := Chirp{
		ID:        newChirp.ID,
		CreatedAt: newChirp.CreatedAt,
		UpdatedAt: newChirp.UpdatedAt,
		Body:      newChirp.Body,
		UserID:    newChirp.UserID,
	}
	RespondWithJson(w, 201, jsonChrip)
}

func cleanText(userText string) string {
	profaneList := []string{"kerfuffle", "sharbert", "fornax"}
	words := strings.Split(userText, " ")
	newWords := make([]string, len(words))
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if slices.Contains(profaneList, loweredWord) {
			newWords[i] = "****"
		} else {
			newWords[i] = word
		}
	}
	return strings.Join(newWords, " ")
}

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	type errBody struct {
		Body string `json:"error"`
	}
	errorChirp := errBody{
		Body: msg,
	}
	errMsg, err := json.Marshal(errorChirp)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(errMsg)
}

func RespondWithJson(w http.ResponseWriter, code int, paylod interface{}) {
	msg, err := json.Marshal(paylod)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(msg)
}
