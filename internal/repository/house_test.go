package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shhesterka04/house-service/internal/dto"
	"github.com/shhesterka04/house-service/internal/repository"
	"github.com/shhesterka04/house-service/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateHouse(t *testing.T) {
	developer := "ABC Corp"
	tests := []struct {
		name      string
		house     *dto.House
		setupMock func(m *mocks.MockDBHouse)
		wantErr   bool
	}{
		{
			name:  "house already exists",
			house: &dto.House{Address: "123 Main St"},
			setupMock: func(m *mocks.MockDBHouse) {
				mockRow := mocks.NewMockRowDBHouse(gomock.NewController(t))
				mockRow.EXPECT().Scan(gomock.Any()).Return(nil).Times(1)
				m.EXPECT().QueryRow(gomock.Any(), "SELECT id FROM house WHERE address = $1", "123 Main St").
					Return(mockRow).Times(1)
			},
			wantErr: true,
		},
		{
			name:  "create house successfully",
			house: &dto.House{Address: "456 Elm St", Year: 2020, Developer: &developer},
			setupMock: func(m *mocks.MockDBHouse) {
				mockRow := mocks.NewMockRowDBHouse(gomock.NewController(t))
				mockRow.EXPECT().Scan(gomock.Any()).Return(pgx.ErrNoRows).Times(1)
				m.EXPECT().QueryRow(gomock.Any(), "SELECT id FROM house WHERE address = $1", "456 Elm St").
					Return(mockRow).Times(1)
				m.EXPECT().Exec(gomock.Any(), "INSERT INTO house (address, year, developer) VALUES ($1, $2, $3)", "456 Elm St", 2020, &developer).
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

			mockDB := mocks.NewMockDBHouse(ctrl)
			tt.setupMock(mockDB)

			repo := repository.NewHouseRepository(mockDB)
			house, err := repo.CreateHouse(context.Background(), tt.house)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.house, house)
			}
		})
	}
}

func TestUpdateHouse(t *testing.T) {
	developer := "ABC Corp"
	developerPtr := &developer
	now := time.Now()
	nowPtr := &now

	tests := []struct {
		name      string
		id        int
		updAt     time.Time
		setupMock func(m *mocks.MockDBHouse)
		wantHouse *dto.House
		wantErr   bool
	}{
		{
			name:  "update house successfully",
			id:    1,
			updAt: now,
			setupMock: func(m *mocks.MockDBHouse) {
				m.EXPECT().Exec(gomock.Any(), "UPDATE house SET updated_at = $1 WHERE id = $2", gomock.Any(), 1).
					Return(pgconn.CommandTag{}, nil).Times(1)
				mockRow := mocks.NewMockRowDBHouse(gomock.NewController(t))
				mockRow.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(dest ...any) error {
						*dest[0].(*int) = 1
						*dest[1].(*string) = "123 Main St"
						*dest[2].(*int) = 2020
						*dest[3].(**string) = developerPtr
						*dest[4].(**time.Time) = nowPtr
						*dest[5].(**time.Time) = nowPtr
						return nil
					}).Times(1)
				m.EXPECT().QueryRow(gomock.Any(), "SELECT id, address, year, developer, created_at, updated_at FROM house WHERE id = $1", 1).
					Return(mockRow).Times(1)
			},
			wantHouse: &dto.House{Id: 1, Address: "123 Main St", Year: 2020, Developer: developerPtr, CreatedAt: nowPtr, UpdateAt: nowPtr},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mocks.NewMockDBHouse(ctrl)
			tt.setupMock(mockDB)

			repo := repository.NewHouseRepository(mockDB)
			house, err := repo.UpdateHouse(context.Background(), tt.id, tt.updAt)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantHouse, house)
			}
		})
	}
}
