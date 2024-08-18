package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shhesterka04/house-service/internal/dto"
	"github.com/shhesterka04/house-service/internal/repository"
	"github.com/shhesterka04/house-service/internal/service"
	"github.com/shhesterka04/house-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.PostRegisterJSONRequestBody
		mockSetup func(m *mocks.MockUserRepo)
		wantErr   bool
	}{
		{
			name: "successful registration",
			req: dto.PostRegisterJSONRequestBody{
				Email:    (*dto.Email)(ptr("test@example.com")),
				Password: ptr("password"),
				UserType: ptr(dto.Client),
			},
			mockSetup: func(m *mocks.MockUserRepo) {
				m.EXPECT().CreateUser(gomock.Any(), gomock.AssignableToTypeOf(repository.User{
					Email:    "test@example.com",
					Password: "",
					Type:     "client",
				})).Return(nil).Times(1)
			},
			wantErr: false,
		},
		{
			name: "invalid user type",
			req: dto.PostRegisterJSONRequestBody{
				Email:    (*dto.Email)(ptr("test@example.com")),
				Password: ptr("password"),
				UserType: (*dto.UserType)(ptr("invalid")),
			},
			mockSetup: func(m *mocks.MockUserRepo) {},
			wantErr:   true,
		},
		{
			name: "error creating user",
			req: dto.PostRegisterJSONRequestBody{
				Email:    (*dto.Email)(ptr("test@example.com")),
				Password: ptr("password"),
				UserType: ptr(dto.Client),
			},
			mockSetup: func(m *mocks.MockUserRepo) {
				m.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errors.New("create user error")).Times(1)
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepo(ctrl)
			tt.mockSetup(mockUserRepo)

			authService := service.NewAuthService(mockUserRepo)
			err := authService.Register(context.Background(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuthService_DummyLogin(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.GetDummyLoginParams
		wantToken string
		wantErr   bool
	}{
		{
			name: "successful dummy login as client",
			req: dto.GetDummyLoginParams{
				UserType: dto.Client,
			},
			wantToken: "client_token",
			wantErr:   false,
		},
		{
			name: "successful dummy login as moderator",
			req: dto.GetDummyLoginParams{
				UserType: dto.Moderator,
			},
			wantToken: gomock.Any().String(),
			wantErr:   false,
		},
		{
			name: "invalid user type",
			req: dto.GetDummyLoginParams{
				UserType: dto.UserType("invalid"),
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			authService := service.NewAuthService(nil)
			_, err := authService.DummyLogin(context.Background(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAuthService_Login(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.LoginRequest
		mockSetup func(m *mocks.MockUserRepo)
		wantToken string
		wantErr   bool
	}{

		{
			name: "invalid email",
			req: dto.LoginRequest{
				Email:    "invalid@example.com",
				Password: "password",
			},
			mockSetup: func(m *mocks.MockUserRepo) {
				m.EXPECT().GetUser(gomock.Any(), "invalid@example.com").Return(repository.User{}, errors.New("invalid email")).Times(1)
			},
			wantToken: "",
			wantErr:   true,
		},
		{
			name: "invalid password",
			req: dto.LoginRequest{
				Email:    "test@example.com",
				Password: "wrongpassword",
			},
			mockSetup: func(m *mocks.MockUserRepo) {
				m.EXPECT().GetUser(gomock.Any(), "test@example.com").Return(repository.User{
					Email:    "test@example.com",
					Password: hashPassword("password"),
					Type:     "client",
				}, nil).Times(1)
			},
			wantToken: "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUserRepo := mocks.NewMockUserRepo(ctrl)
			tt.mockSetup(mockUserRepo)

			authService := service.NewAuthService(mockUserRepo)
			token, err := authService.Login(context.Background(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantToken, token)
			}
		})
	}
}

func ptr[T interface{}](v T) *T {
	return &v
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
