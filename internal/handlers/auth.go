package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/shhesterka04/house-service/internal/repository"
	"github.com/shhesterka04/house-service/pkg/logger"
)

type LoginRequest struct {
	UserType string `json:"user_type"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type DummyLoginRequest struct {
	UserType string `json:"user_type"`
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

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	user := repository.User{
		Email:    req.Email,
		Password: hashedPassword,
		Type:     req.UserType,
	}

	if err = s.userRepo.CreateUser(r.Context(), user); err != nil {
		http.Error(w, "create user", http.StatusBadRequest)
		return
	}

	logger.Infof(r.Context(), "User %s registered", req.Email)
	w.WriteHeader(http.StatusCreated)
}

func (s *AuthService) Login(w http.ResponseWriter, r *http.Request) {
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

	if !checkPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid password", http.StatusBadRequest)
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
	logger.Infof(r.Context(), "User %s logged in", req.Email)
	json.NewEncoder(w).Encode(res)
}
