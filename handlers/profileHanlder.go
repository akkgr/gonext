package handlers

import (
	"net/http"
)

func (h *Handler) getProfile(w http.ResponseWriter, r *http.Request) {
	returnJSON(http.StatusOK, h.Claims["sub"].(string), w)
}
