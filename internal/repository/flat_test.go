package repository_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/shhesterka04/house-service/internal/dto"
	"github.com/shhesterka04/house-service/internal/repository"
	"github.com/shhesterka04/house-service/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateFlat(t *testing.T) {
	tests := []struct {
		name      string
		flat      *dto.DtoFlat
		setupMock func(m *mocks.MockDBFlat)
		wantErr   bool
	}{
		{
			name: "flat already exists",
			flat: &dto.DtoFlat{HouseId: 1, Number: 101},
			setupMock: func(m *mocks.MockDBFlat) {
				mockRow := mocks.NewMockRowDBFlat(gomock.NewController(t))
				mockRow.EXPECT().Scan(gomock.Any()).Return(nil).Times(1)
				m.EXPECT().QueryRow(gomock.Any(), "SELECT id FROM flats WHERE house_id = $1 AND number = $2", 1, 101).
					Return(mockRow).Times(1)
			},
			wantErr: true,
		},
		{
			name: "create flat successfully",
			flat: &dto.DtoFlat{HouseId: 1, Status: "created", Number: 102, Rooms: 3, Price: 100000},
			setupMock: func(m *mocks.MockDBFlat) {
				mockRow := mocks.NewMockRowDBFlat(gomock.NewController(t))
				mockRow.EXPECT().Scan(gomock.Any()).Return(pgx.ErrNoRows).Times(1)
				m.EXPECT().QueryRow(gomock.Any(), "SELECT id FROM flats WHERE house_id = $1 AND number = $2", 1, 102).
					Return(mockRow).Times(1)
				m.EXPECT().Exec(gomock.Any(), "INSERT INTO flats (house_id, status, number, rooms, price) VALUES ($1, $2, $3, $4, $5)", 1, "created", 102, 3, 100000).
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

			mockDB := mocks.NewMockDBFlat(ctrl)
			tt.setupMock(mockDB)

			repo := repository.NewFlatRepository(mockDB)
			flat, err := repo.CreateFlat(context.Background(), tt.flat)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.flat, flat)
			}
		})
	}
}

func TestGetFlatByID(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		setupMock func(m *mocks.MockDBFlat)
		wantFlat  *dto.DtoFlat
		wantErr   bool
	}{
		{
			name: "flat not found",
			id:   1,
			setupMock: func(m *mocks.MockDBFlat) {
				mockRow := mocks.NewMockRowDBFlat(gomock.NewController(t))
				mockRow.EXPECT().Scan(gomock.Any()).Return(pgx.ErrNoRows).Times(1)
				m.EXPECT().QueryRow(gomock.Any(), "SELECT id, house_id, status, number, rooms, price FROM flats WHERE id = $1", 1).
					Return(mockRow).Times(1)
			},
			wantFlat: nil,
			wantErr:  true,
		},
		{
			name: "get flat successfully",
			id:   2,
			setupMock: func(m *mocks.MockDBFlat) {
				mockRow := mocks.NewMockRowDBFlat(gomock.NewController(t))
				mockRow.EXPECT().Scan(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(
					func(dest ...any) error {
						*dest[0].(*int) = 2
						*dest[1].(*int) = 1
						*dest[2].(*string) = "created"
						*dest[3].(*int) = 102
						*dest[4].(*int) = 3
						*dest[5].(*int) = 100000
						return nil
					}).Times(1)
				m.EXPECT().QueryRow(gomock.Any(), "SELECT id, house_id, status, number, rooms, price FROM flats WHERE id = $1", 2).
					Return(mockRow).Times(1)
			},
			wantFlat: &dto.DtoFlat{Id: 2, HouseId: 1, Status: "created", Number: 102, Rooms: 3, Price: 100000},
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockDB := mocks.NewMockDBFlat(ctrl)
			tt.setupMock(mockDB)

			repo := repository.NewFlatRepository(mockDB)
			flat, err := repo.GetFlatByID(context.Background(), tt.id)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantFlat, flat)
			}
		})
	}
}

// TODO
func TestGetFlatByHouseID(t *testing.T) {
}
