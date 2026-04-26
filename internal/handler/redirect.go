package handler

import (
	"net/http"
	"regexp"

	"github.com/Wild-sergunys/shrtic/internal/service"
)

var shortCodeRegex = regexp.MustCompile(`^/[a-zA-Z0-9]{7}$`)

type RedirectHandler struct {
	linkService *service.LinkService
}

func NewRedirectHandler(linkService *service.LinkService) *RedirectHandler {
	return &RedirectHandler{linkService: linkService}
}

func (h *RedirectHandler) RedirectToLongURL(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Path

	if !shortCodeRegex.MatchString(code) {
		writeError(w, http.StatusNotFound, "not_found", "Ссылка не найдена", nil)
		return
	}

	longURL, err := h.linkService.GetLongURL(r.Context(), code)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", "Ссылка не найдена", nil)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}
