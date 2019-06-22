package log

import (
	"net/http"
	"time"
)

func (h *Handler) log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("request processed in %s from %s\n", time.Now().Sub(startTime), r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
