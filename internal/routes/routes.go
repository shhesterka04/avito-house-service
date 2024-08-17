package routes

import (
	"net/http"

	"github.com/shhesterka04/house-service/internal/dto"
	"github.com/shhesterka04/house-service/internal/handlers"
	"github.com/shhesterka04/house-service/internal/middleware"
)

func NewRouter(authHandlers *handlers.AuthHandlers, houseHandlers *handlers.HouseHandler, flatHandlers *handlers.FlatHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /dummyLogin", authHandlers.DummyLogin)
	mux.HandleFunc("POST /login", authHandlers.Login)
	mux.HandleFunc("POST /register", authHandlers.Register)

	protectedRoutes := http.NewServeMux()
	protectedRoutes.Handle("POST /house/create", middleware.AuthMiddleware(dto.Moderator)(http.HandlerFunc(houseHandlers.CreateHouse)))
	protectedRoutes.Handle("GET /house/{id}", middleware.AuthMiddleware(dto.Client)(http.HandlerFunc(flatHandlers.GetFlatsByHouseID)))
	protectedRoutes.Handle("POST /house/{id}/subscribe", middleware.AuthMiddleware(dto.Client)(http.HandlerFunc(houseHandlers.SubscribeToHouse)))
	protectedRoutes.Handle("POST /flat/create", middleware.AuthMiddleware(dto.Client)(http.HandlerFunc(flatHandlers.CreateFlat)))
	protectedRoutes.Handle("POST /flat/update", middleware.AuthMiddleware(dto.Moderator)(http.HandlerFunc(flatHandlers.UpdateFlat)))

	mux.Handle("/", protectedRoutes)
	return mux
}
