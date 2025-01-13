package handlers

import "net/http"

func (h *UserHandler) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		h.cfg.FileserverHits.Add(1)
		next.ServeHTTP(w, req)
	})
}
