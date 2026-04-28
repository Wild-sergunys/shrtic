package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/Wild-sergunys/shrtic/internal/middleware"
	"github.com/Wild-sergunys/shrtic/internal/model"
	"github.com/Wild-sergunys/shrtic/internal/service"
)

var loginRateLimiter *middleware.RateLimiter

func SetLoginRateLimiter(rl *middleware.RateLimiter) {
	loginRateLimiter = rl
}

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req model.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Неверный формат запроса", nil)
		return
	}

	if req.Login == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "Логин и пароль обязательны", nil)
		return
	}

	if len(req.Login) < 3 || len(req.Login) > 50 {
		writeError(w, http.StatusBadRequest, "validation_error", "Логин должен быть от 3 до 50 символов", nil)
		return
	}

	user, err := h.authService.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		if strings.Contains(err.Error(), "уже существует") {
			writeError(w, http.StatusConflict, "conflict", err.Error(), nil)
			return
		}
		writeError(w, http.StatusBadRequest, "validation_error", err.Error(), nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(model.IdResponse{ID: user.ID})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req model.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "validation_error", "Неверный формат запроса", nil)
		return
	}

	if req.Login == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "validation_error", "Логин и пароль обязательны", nil)
		return
	}

	token, err := h.authService.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		if err == service.ErrInvalidCredentials {
			writeError(w, http.StatusUnauthorized, "invalid_credentials", err.Error(), nil)
			return
		}
		writeError(w, http.StatusInternalServerError, "internal_error", "Внутренняя ошибка сервера", nil)
		return
	}

	if loginRateLimiter != nil {
		loginRateLimiter.Reset(middleware.GetIP(r))
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.LoginResponse{
		Token: token,
		Role:  "user",
	})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(model.UserIDKey).(int64)
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized", "Не авторизован", nil)
		return
	}

	user, err := h.authService.GetUser(r.Context(), userID)
	if err != nil || user == nil {
		writeError(w, http.StatusInternalServerError, "internal_error", "Ошибка сервера", nil)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(model.MessageResponse{Message: "Выход выполнен успешно"})
}
