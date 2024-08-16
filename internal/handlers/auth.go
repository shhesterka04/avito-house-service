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
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.authService.Register(r.Context(), req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	logger.Infof(r.Context(), "User %h registered", req.Email)
	w.WriteHeader(http.StatusCreated)
}

func (h *AuthHandlers) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req dto.GetDummyLoginParams
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	token, err := h.authService.DummyLogin(r.Context(), req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

// TODO: разобраться с openAPI
func (h *AuthHandlers) Login(w http.ResponseWriter, r *http.Request) {
	//var req dto.PostLoginJSONRequestBody
	//if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
	//	http.Error(w, "Invalid request payload", http.StatusBadRequest)
	//	return
	//}
	//
	//token, err := h.authService.Login(r.Context(), req)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//
	//res := dto.LoginResponse{Token: token}
	//w.Header().Set("Content-Type", "application/json")
	//logger.Infof(r.Context(), "User %h logged in", req.Email)
	//json.NewEncoder(w).Encode(res)
}
