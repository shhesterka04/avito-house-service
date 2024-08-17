package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/shhesterka04/house-service/internal/dto"
	"github.com/shhesterka04/house-service/internal/service"
	"github.com/shhesterka04/house-service/pkg/logger"
)

type AuthHandlers struct {
	authService *service.AuthService
}

func NewAuthHandlers(authService *service.AuthService) *AuthHandlers {
	return &AuthHandlers{authService: authService}
}

func (h *AuthHandlers) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.PostRegisterJSONRequestBody
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Errorf(r.Context(), "Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.authService.Register(r.Context(), req); err != nil {
		logger.Errorf(r.Context(), "Error registering user: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Debugf(r.Context(), "User %h registered", *req.Email)
	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandlers) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.GetDummyLoginParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Errorf(r.Context(), "Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.authService.DummyLogin(r.Context(), req)
	if err != nil {
		logger.Errorf(r.Context(), "Error dummy logging in: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Errorf(r.Context(), "Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.authService.Login(r.Context(), req)
	if err != nil {
		logger.Errorf(r.Context(), "Error logging in: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	logger.Debugf(r.Context(), "User %h logged in", req.Email)
	json.NewEncoder(w).Encode(token)
}
