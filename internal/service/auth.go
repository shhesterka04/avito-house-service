package service

import (
	"context"

	"github.com/pkg/errors"
	"github.com/shhesterka04/house-service/internal/dto"
	"github.com/shhesterka04/house-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUserType = errors.New("invalid user type")
	ErrHashPassword    = errors.New("hash password")
)

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

func (s *AuthService) Register(ctx context.Context, req dto.PostRegisterJSONRequestBody) error {
	if *req.UserType != dto.Client && *req.UserType != dto.Moderator {
		return ErrInvalidUserType
	}

	hashedPassword, err := hashPassword(*req.Password)
	if err != nil {
		return ErrHashPassword
	}

	user := repository.User{
		Email:    string(*req.Email),
		Password: hashedPassword,
		Type:     string(*req.UserType),
	}

	if err = s.userRepo.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *AuthService) DummyLogin(ctx context.Context, req dto.GetDummyLoginParams) (string, error) {
	var token string
	//TODO: сделать норм токены
	switch req.UserType {
	case dto.Client:
		token = "client_token"
	case dto.Moderator:
		token = "moderator_token"
	default:
		return "", ErrInvalidUserType
	}

	return token, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	user, err := s.userRepo.GetUser(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid email")
	}

	if !checkPasswordHash(req.Password, user.Password) {
		return "", errors.New("invalid password")
	}

	var token string
	switch user.Type {
	case "client":
		token = "client_token"
	case "moderator":
		token = "moderator_token"
	default:
		return "", errors.New("invalid user type")
	}

	return token, nil
}
