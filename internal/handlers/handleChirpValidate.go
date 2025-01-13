package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func (h *UserHandler) HandleChirpValidate(w http.ResponseWriter, req *http.Request) {
	type chirp struct {
		Body string `json:"body"`
	}

	type validBody struct {
		Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(req.Body)
	newChirp := chirp{}
	err := decoder.Decode(&newChirp)
	if err != nil {
		errMsg := "Something went wrong"
		RespondWithError(w, 400, errMsg)
		return
	}

	if len(newChirp.Body) > 140 {
		errMsg := "Chirp is too long"
		RespondWithError(w, 400, errMsg)
		return
	}

	cleanedBody := cleanText(newChirp.Body)

	validChirp := validBody{
		Body: cleanedBody,
	}
	RespondWithJson(w, 200, validChirp)
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
