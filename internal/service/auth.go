//go:generate mockgen -source ./auth.go -destination=./mocks/auth.go -package=mocks
package service

import (
	"context"
	"regexp"

	"github.com/pkg/errors"
	"github.com/shhesterka04/house-service/internal/dto"
	"github.com/shhesterka04/house-service/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidUserType = errors.New("invalid user type")

	ErrorInvalidLogin = errors.New("invalid login")
	ErrInValidEmail   = errors.New("invalid email")

	re = regexp.MustCompile(emailRegex)
)

const (
	emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
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

func isValidEmail(email string) bool {
	return re.MatchString(email)
}

func (s *AuthService) Register(ctx context.Context, req dto.PostRegisterJSONRequestBody) error {
	if *req.UserType != dto.Client && *req.UserType != dto.Moderator {
		return ErrInvalidUserType
	}

	if req.Email == nil || *req.Email == "" || !isValidEmail(string(*req.Email)) {
		return ErrInValidEmail
	}

	if req.Password == nil || *req.Password == "" {
		return ErrorInvalidLogin
	}

	hashedPassword, err := hashPassword(*req.Password)
	if err != nil {
		return ErrorInvalidLogin
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
	var userType string
	switch req.UserType {
	case dto.Client:
		userType = "client"
	case dto.Moderator:
		userType = "moderator"
	default:
		return "", ErrInvalidUserType
	}

	token, err := GenerateJWT(userType)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (string, error) {
	user, err := s.userRepo.GetUser(ctx, req.Email)
	if err != nil {
		return "", ErrorInvalidLogin
	}

	if !checkPasswordHash(req.Password, user.Password) {
		return "", ErrorInvalidLogin
	}

	token, err := GenerateJWT(user.Type)
	if err != nil {
		return "", err
	}

	return token, nil
}
