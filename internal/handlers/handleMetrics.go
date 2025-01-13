package handlers

import (
	"fmt"
	"net/http"
)

func (h *UserHandler) MetricsHandler(w http.ResponseWriter, req *http.Request) {
	const metricsHtml = `
	<html>
	  <body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	  </body>
	</html>
	`
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, metricsHtml, h.cfg.FileserverHits.Load())
}
