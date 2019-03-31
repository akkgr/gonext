package handlers

import "net/http"

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	returnText(http.StatusOK, h.claims["sub"].(string), w)
}
