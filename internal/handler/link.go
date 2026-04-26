package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/Wild-sergunys/shrtic/internal/model"
	"github.com/Wild-sergunys/shrtic/internal/service"
)

type LinkHandler struct {
	linkService *service.LinkService
}

func NewLinkHandler(linkService *service.LinkService) *LinkHandler {
	return &LinkHandler{linkService: linkService}
}

func (h *LinkHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	var req model.CreateLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Неверный формат запроса", nil)
		return
	}

	if req.URL == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "URL обязателен", nil)
		return
	}

	var userID *int64
	if uid, ok := r.Context().Value(model.UserIDKey).(int64); ok {
		userID = &uid
	}

	link, err := h.linkService.CreateShortLink(r.Context(), req.URL, userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(link)
}

func (h *LinkHandler) GetLinks(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(model.UserIDKey).(int64)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Не авторизован", nil)
		return
	}

	search := r.URL.Query().Get("search")
	links, err := h.linkService.GetLinks(r.Context(), userID, search)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Ошибка сервера", nil)
		return
	}

	if links == nil {
		links = []model.Link{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}

func (h *LinkHandler) DeleteLink(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(model.UserIDKey).(int64)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Не авторизован", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Неверный ID", nil)
		return
	}

	if err := h.linkService.DeleteLink(r.Context(), userID, id); err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			writeError(w, http.StatusNotFound, "not_found", err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "нет прав") {
			writeError(w, http.StatusForbidden, "forbidden", err.Error(), nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Ошибка сервера", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.MessageResponse{Message: "Ссылка удалена"})
}

func (h *LinkHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(model.UserIDKey).(int64)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Не авторизован", nil)
		return
	}

	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Неверный ID", nil)
		return
	}

	stats, err := h.linkService.GetStats(r.Context(), userID, id)
	if err != nil {
		if strings.Contains(err.Error(), "не найдена") {
			writeError(w, http.StatusNotFound, "not_found", err.Error(), nil)
			return
		}
		if strings.Contains(err.Error(), "нет прав") {
			writeError(w, http.StatusForbidden, "forbidden", err.Error(), nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Ошибка сервера", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}
