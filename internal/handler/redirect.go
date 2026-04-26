package handler

import (
	"context"
	"net"
	"net/http"
	"regexp"
	"strings"

	"github.com/Wild-sergunys/shrtic/internal/service"
)

var shortCodeRegex = regexp.MustCompile(`^/r/[a-zA-Z0-9]{7}$`)

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

	go func() {
		link, err := h.linkService.GetLinkByCode(context.Background(), code)
		if err == nil && link != nil {
			h.linkService.RecordClick(context.Background(), link.ID)
			clientIP := getIP(r)
			h.linkService.SaveClickStats(context.Background(), link.ID, r.UserAgent(), r.Referer(), clientIP)
		}
	}()

	http.Redirect(w, r, longURL, http.StatusFound)
}

func getIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ip := strings.Split(xff, ",")[0]
		return strings.TrimSpace(ip)
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return strings.TrimSpace(xri)
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
