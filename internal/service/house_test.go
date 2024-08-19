//go:build unit
// +build unit

package service_test

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/shhesterka04/house-service/internal/dto"
	"github.com/shhesterka04/house-service/internal/service"
	"github.com/shhesterka04/house-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func Test_ServiceCreateHouse(t *testing.T) {
	developer := "Developer Inc."

	tests := []struct {
		name      string
		req       dto.PostHouseCreateJSONRequestBody
		mockSetup func(m *mocks.MockHouseRepo)
		wantHouse *dto.House
		wantErr   bool
	}{
		{
			name: "successful creation",
			req: dto.PostHouseCreateJSONRequestBody{
				Address:   "123 Main St",
				Year:      2020,
				Developer: &developer,
			},
			mockSetup: func(m *mocks.MockHouseRepo) {
				m.EXPECT().CreateHouse(gomock.Any(), &dto.House{
					Address:   "123 Main St",
					Year:      2020,
					Developer: &developer,
				}).Return(&dto.House{
					Address:   "123 Main St",
					Year:      2020,
					Developer: &developer,
				}, nil).Times(1)
			},
			wantHouse: &dto.House{
				Address:   "123 Main St",
				Year:      2020,
				Developer: &developer,
			},
			wantErr: false,
		},
		{
			name: "house already exists",
			req: dto.PostHouseCreateJSONRequestBody{
				Address:   "123 Main St",
				Year:      2020,
				Developer: &developer,
			},
			mockSetup: func(m *mocks.MockHouseRepo) {
				m.EXPECT().CreateHouse(gomock.Any(), &dto.House{
					Address:   "123 Main St",
					Year:      2020,
					Developer: &developer,
				}).Return(nil, errors.New("house already exists")).Times(1)
			},
			wantHouse: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockHouseRepo := mocks.NewMockHouseRepo(ctrl)
			tt.mockSetup(mockHouseRepo)

			houseService := service.NewHouseService(mockHouseRepo)
			house, err := houseService.CreateHouse(context.Background(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantHouse, house)
			}
		})
	}
}
