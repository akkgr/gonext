package handlers

import (
	"net/http"
)

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	returnJSON(http.StatusOK, "api status ok", w)
}
