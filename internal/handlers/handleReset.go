package handlers

import "net/http"

func (h *UserHandler) ResetHandler(w http.ResponseWriter, req *http.Request) {

	if h.cfg.Platform != "dev" {
		errMsg := "Error, action is forbidden"
		RespondWithError(w, http.StatusForbidden, errMsg)
		return
	}
	h.cfg.FileserverHits.Store(0)
	err := h.cfg.DB.DeleteUsers(req.Context())
	if err != nil {
		errMsg := "Error, could not delete users"
		RespondWithError(w, 400, errMsg)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Users deleted successfully"))
}
