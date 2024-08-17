package repository_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shhesterka04/house-service/internal/repository"
	"github.com/shhesterka04/house-service/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name      string
		user      repository.User
		setupMock func(m *mock_db.MockDBUser)
		wantErr   bool
	}{
		{
			name: "user already exists",
			user: repository.User{Email: "test@example.com", Password: "password", Type: "client"},
			setupMock: func(m *mock_db.MockDBUser) {
				mockRow := mock_db.NewMockRow(gomock.NewController(t))
				mockRow.EXPECT().Scan(gomock.Any()).Return(nil).Times(1)
				m.EXPECT().QueryRow(gomock.Any(), "SELECT id FROM users WHERE email = $1", "test@example.com").
					Return(mockRow).Times(1)
			},
			wantErr: true,
		},
		{
			name: "create user successfully",
			user: repository.User{Email: "newuser@example.com", Password: "password", Type: "client"},
			setupMock: func(m *mock_db.MockDBUser) {
				mockRow := mock_db.NewMockRow(gomock.NewController(t))
				mockRow.EXPECT().Scan(gomock.Any()).Return(pgx.ErrNoRows).Times(1)
				m.EXPECT().QueryRow(gomock.Any(), "SELECT id FROM users WHERE email = $1", "newuser@example.com").
					Return(mockRow).Times(1)
				m.EXPECT().Exec(gomock.Any(), "INSERT INTO users (email, password, type) VALUES ($1, $2, $3)", "newuser@example.com", "password", "client").
					Return(pgconn.CommandTag{}, nil).Times(1)
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_db.NewMockDBUser(ctrl)
			tt.setupMock(mockDB)

			repo := repository.NewUserRepository(mockDB)
			err := repo.CreateUser(context.Background(), tt.user)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetUser(t *testing.T) {
	tests := []struct {
		name      string
		email     string
		setupMock func(m *mock_db.MockDBUser)
		wantUser  repository.User
		wantErr   bool
	}{
		{
			name:  "user not found",
			email: "notfound@example.com",
			setupMock: func(m *mock_db.MockDBUser) {
				mockRow := mock_db.NewMockRow(gomock.NewController(t))
				mockRow.EXPECT().Scan(gomock.Any()).Return(pgx.ErrNoRows).Times(1)
				m.EXPECT().QueryRow(gomock.Any(), "SELECT * FROM users WHERE email = $1", "notfound@example.com").
					Return(mockRow).Times(1)
			},
			wantUser: repository.User{},
			wantErr:  true,
		},
		{
			name:  "get user successfully",
			email: "test@example.com",
			setupMock: func(m *mock_db.MockDBUser) {
				mockRow := mock_db.NewMockRow(gomock.NewController(t))
				mockRow.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(dest ...any) error {
						*dest[0].(*string) = "uuid"
						*dest[1].(*string) = "test@example.com"
						*dest[2].(*string) = "password"
						*dest[3].(*string) = "client"
						return nil
					}).Times(1)
				m.EXPECT().QueryRow(gomock.Any(), "SELECT * FROM users WHERE email = $1", "test@example.com").
					Return(mockRow).Times(1)
			},
			wantUser: repository.User{UUID: "uuid", Email: "test@example.com", Password: "password", Type: "client"},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mock_db.NewMockDBUser(ctrl)
			tt.setupMock(mockDB)

			repo := repository.NewUserRepository(mockDB)
			user, err := repo.GetUser(context.Background(), tt.email)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantUser, user)
			}
		})
	}
}
