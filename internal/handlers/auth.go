package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/shhesterka04/house-service/internal/repository"
)

type LoginRequest struct {
	UserType string `json:"user_type"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UserRepo interface {
	CreateUser(ctx context.Context, user repository.User) error
	GetUser(ctx context.Context, email string) (repository.User, error)
}

type AuthService struct {
	userRepo UserRepo
}

func NewAuthService(userRepo UserRepo) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) DummyLogin(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := s.userRepo.GetUser(r.Context(), req.Email)
	if err != nil {
		http.Error(w, "Invalid email", http.StatusBadRequest)
		return
	}

	var token string

	switch user.Type {
	case "client":
		token = "client_token"
	case "moderator":
		token = "moderator_token"
	default:
		http.Error(w, "Invalid user type", http.StatusBadRequest)
		return
	}

	res := LoginResponse{Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (s *AuthService) Register(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.UserType != "client" && req.UserType != "moderator" {
		http.Error(w, "Invalid user type", http.StatusBadRequest)
		return
	}

	user := repository.User{
		Email:    req.Email,
		Password: req.Password,
		Type:     req.UserType,
	}

	if err := s.userRepo.CreateUser(r.Context(), user); err != nil {
		http.Error(w, "User already exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {}
