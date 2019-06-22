package app

import (
	"net/http"
)

func (h *Handler) root(w http.ResponseWriter, r *http.Request) {
	returnText(http.StatusOK, "api status ok", w)
}
