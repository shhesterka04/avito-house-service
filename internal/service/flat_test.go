//go:build unit
// +build unit

package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/shhesterka04/house-service/internal/dto"
	"github.com/shhesterka04/house-service/internal/service"
	"github.com/shhesterka04/house-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestService_CreateFlat(t *testing.T) {
	tests := []struct {
		name      string
		req       dto.CreateFlatRequest
		mockSetup func(m *mocks.MockFlatRepo, h *mocks.MockHouseFlatRepo)
		wantFlat  *dto.DtoFlat
		wantErr   bool
	}{
		{
			name: "successful creation",
			req: dto.CreateFlatRequest{
				HouseID: 1,
				Number:  101,
				Rooms:   3,
				Price:   100000,
			},
			mockSetup: func(m *mocks.MockFlatRepo, h *mocks.MockHouseFlatRepo) {
				m.EXPECT().CreateFlat(gomock.Any(), &dto.DtoFlat{
					HouseID: 1,
					Number:  101,
					Rooms:   3,
					Price:   100000,
					Status:  string(dto.Created),
				}).Return(&dto.DtoFlat{
					HouseID: 1,
					Number:  101,
					Rooms:   3,
					Price:   100000,
					Status:  string(dto.Created),
				}, nil).Times(1)
				h.EXPECT().UpdateHouse(gomock.Any(), 1, gomock.Any()).Return(nil, nil).Times(1)
			},
			wantFlat: &dto.DtoFlat{
				HouseID: 1,
				Number:  101,
				Rooms:   3,
				Price:   100000,
				Status:  string(dto.Created),
			},
			wantErr: false,
		},
		{
			name: "error creating flat",
			req: dto.CreateFlatRequest{
				HouseID: 1,
				Number:  101,
				Rooms:   3,
				Price:   100000,
			},
			mockSetup: func(m *mocks.MockFlatRepo, h *mocks.MockHouseFlatRepo) {
				m.EXPECT().CreateFlat(gomock.Any(), gomock.Any()).Return(nil, errors.New("create flat error")).Times(1)
			},
			wantFlat: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFlatRepo := mocks.NewMockFlatRepo(ctrl)
			mockHouseFlatRepo := mocks.NewMockHouseFlatRepo(ctrl)
			tt.mockSetup(mockFlatRepo, mockHouseFlatRepo)

			flatService := service.NewFlatService(mockFlatRepo, mockHouseFlatRepo)
			flat, err := flatService.CreateFlat(context.Background(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantFlat, flat)
			}
		})
	}
}

func TestFlatService_UpdateFlat(t *testing.T) {
	validStatus := dto.Approved
	invalidStatus := dto.Status("invalid")

	tests := []struct {
		name      string
		req       dto.PostFlatUpdateJSONRequestBody
		mockSetup func(m *mocks.MockFlatRepo, h *mocks.MockHouseFlatRepo)
		wantFlat  *dto.DtoFlat
		wantErr   bool
	}{
		{
			name: "successful update",
			req: dto.PostFlatUpdateJSONRequestBody{
				Id:     1,
				Status: &validStatus,
			},
			mockSetup: func(m *mocks.MockFlatRepo, h *mocks.MockHouseFlatRepo) {
				m.EXPECT().GetFlatByID(gomock.Any(), 1).Return(&dto.DtoFlat{
					ID:     1,
					Status: string(dto.Created),
				}, nil).Times(1)
				m.EXPECT().UpdateFlat(gomock.Any(), &dto.DtoFlat{
					ID:     1,
					Status: string(validStatus),
				}).Return(&dto.DtoFlat{
					ID:     1,
					Status: string(validStatus),
				}, nil).Times(1)
				h.EXPECT().UpdateHouse(gomock.Any(), gomock.Eq(0), gomock.Any()).Return(nil, nil).Times(1)
			},
			wantFlat: &dto.DtoFlat{
				ID:     1,
				Status: string(validStatus),
			},
			wantErr: false,
		},
		{
			name: "invalid status",
			req: dto.PostFlatUpdateJSONRequestBody{
				Id:     1,
				Status: &invalidStatus,
			},
			mockSetup: func(m *mocks.MockFlatRepo, h *mocks.MockHouseFlatRepo) {},
			wantFlat:  nil,
			wantErr:   true,
		},
		{
			name: "flat not found",
			req: dto.PostFlatUpdateJSONRequestBody{
				Id:     1,
				Status: &validStatus,
			},
			mockSetup: func(m *mocks.MockFlatRepo, h *mocks.MockHouseFlatRepo) {
				m.EXPECT().GetFlatByID(gomock.Any(), 1).Return(nil, errors.New("DtoFlat not found")).Times(1)
			},
			wantFlat: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFlatRepo := mocks.NewMockFlatRepo(ctrl)
			mockHouseFlatRepo := mocks.NewMockHouseFlatRepo(ctrl)
			tt.mockSetup(mockFlatRepo, mockHouseFlatRepo)

			flatService := service.NewFlatService(mockFlatRepo, mockHouseFlatRepo)
			flat, err := flatService.UpdateFlat(context.Background(), tt.req)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantFlat, flat)
			}
		})
	}
}

func TestFlatService_GetFlatsByHouseID(t *testing.T) {
	tests := []struct {
		name      string
		houseID   string
		token     string
		mockSetup func(m *mocks.MockFlatRepo)
		wantFlats []*dto.DtoFlat
		wantErr   bool
	}{
		{
			name:    "successful retrieval",
			houseID: "1",
			token:   "",
			mockSetup: func(m *mocks.MockFlatRepo) {
				m.EXPECT().GetFlatByHouseID(gomock.Any(), 1, "client").Return([]*dto.DtoFlat{
					{
						ID:      1,
						HouseID: 1,
						Number:  101,
						Rooms:   3,
						Price:   100000,
						Status:  string(dto.Created),
					},
				}, nil).Times(1)
			},
			wantFlats: []*dto.DtoFlat{
				{
					ID:      1,
					HouseID: 1,
					Number:  101,
					Rooms:   3,
					Price:   100000,
					Status:  string(dto.Created),
				},
			},
			wantErr: false,
		},
		{
			name:      "invalid house ID",
			houseID:   "invalid",
			token:     "",
			mockSetup: func(m *mocks.MockFlatRepo) {},
			wantFlats: nil,
			wantErr:   true,
		},
		{
			name:      "invalid token",
			houseID:   "1",
			token:     "invalid-token",
			mockSetup: func(m *mocks.MockFlatRepo) {},
			wantFlats: nil,
			wantErr:   true,
		},
		{
			name:    "error retrieving flats",
			houseID: "1",
			token:   "",
			mockSetup: func(m *mocks.MockFlatRepo) {
				m.EXPECT().GetFlatByHouseID(gomock.Any(), 1, "client").Return(nil, errors.New("error retrieving flats")).Times(1)
			},
			wantFlats: nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockFlatRepo := mocks.NewMockFlatRepo(ctrl)
			tt.mockSetup(mockFlatRepo)

			if tt.token == "" {
				token, err := service.GenerateJWT("client")
				require.NoError(t, err)
				tt.token = token
			}

			flatService := service.NewFlatService(mockFlatRepo, nil)
			flats, err := flatService.GetFlatsByHouseID(context.Background(), tt.houseID, tt.token)

			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantFlats, flats)
			}
		})
	}
}
