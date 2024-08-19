//go:build integration
// +build integration

package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/shhesterka04/house-service/internal/config"
	"github.com/shhesterka04/house-service/internal/handlers"
	"github.com/shhesterka04/house-service/internal/repository"
	"github.com/shhesterka04/house-service/internal/routes"
	"github.com/shhesterka04/house-service/internal/service"
	"github.com/shhesterka04/house-service/pkg/db"
	"github.com/shhesterka04/house-service/pkg/logger"
	"github.com/stretchr/testify/assert"
)

const (
	migrationDir  = ".././migrations"
	testConfigDir = "testConfig.env"
)

const (
	levelConfig = "debug"
	development = true
)

func TestIntegration(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadConfig(testConfigDir)
	assert.NoError(t, err)

	logger.Init(
		logger.Config{
			Level:       levelConfig,
			Development: development,
		})

	pgClient := db.NewClient(
		cfg.DBName,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
	)

	dbConn, err := pgClient.Connect(ctx)
	assert.NoError(t, err)

	err = pgClient.ResetMigrations(migrationDir)
	assert.NoError(t, err)
	err = pgClient.MigrateUp(migrationDir)
	assert.NoError(t, err)
	defer pgClient.Close()

	userRepo := repository.NewUserRepository(dbConn.Cluster)
	authService := service.NewAuthService(userRepo)
	authHandlers := handlers.NewAuthHandlers(authService)

	houseRepo := repository.NewHouseRepository(dbConn.Cluster)
	houseService := service.NewHouseService(houseRepo)
	houseHandlers := handlers.NewHouseHandler(houseService)

	flatRepo := repository.NewFlatRepository(dbConn.Cluster)
	flatService := service.NewFlatService(flatRepo, houseRepo)
	flatHandlers := handlers.NewFlatHandler(flatService)

	mux := routes.NewRouter(authHandlers, houseHandlers, flatHandlers)

	server := &http.Server{
		Addr:    cfg.HostAddr,
		Handler: mux,
	}

	go func() {
		logger.Infof(ctx, "starting server on %s", cfg.HostAddr)
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			t.Fatalf("listen: %s\n", err)
		}
	}()

	time.Sleep(1 * time.Second)

	client := &http.Client{}

	// Step 1: Register user
	registerPayload := map[string]string{
		"email":     "abaac@lmao.com",
		"password":  "qwerty",
		"user_type": "moderator",
	}

	registerBody, _ := json.Marshal(registerPayload)
	resp, err := client.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(registerBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Step 2: Login user
	loginPayload := map[string]string{
		"email":    "abaac@lmao.com",
		"password": "qwerty",
	}
	loginBody, _ := json.Marshal(loginPayload)
	resp, err = client.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(loginBody))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	tokenBytes, _ := io.ReadAll(resp.Body)
	token := string(tokenBytes)
	token = token[1 : len(token)-2]
	fmt.Println(token)

	// Step 3: Create house
	housePayload := map[string]interface{}{
		"address":   "123 Main St",
		"year":      2020,
		"developer": "Developer Inc.",
	}
	houseBody, _ := json.Marshal(housePayload)
	req, _ := http.NewRequest("POST", "http://localhost:8080/house/create", bytes.NewBuffer(houseBody))
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	fmt.Println(resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Get house ID
	respHouse, _ := io.ReadAll(resp.Body)
	var house map[string]interface{}
	json.Unmarshal(respHouse, &house)
	houseIDd := house["id"].(float64)
	houseID := int(houseIDd)

	// Step 4: Create flats
	for i := 1; i <= 3; i++ {
		flatPayload := map[string]interface{}{
			"house_id": houseID,
			"number":   i,
			"rooms":    2,
			"price":    100000 * i,
		}
		flatBody, _ := json.Marshal(flatPayload)
		req, _ = http.NewRequest("POST", "http://localhost:8080/flat/create", bytes.NewBuffer(flatBody))
		req.Header.Set("Authorization", "Bearer "+token)
		resp, err = client.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	}
	respFlat, _ := io.ReadAll(resp.Body)
	var flat map[string]interface{}
	json.Unmarshal(respFlat, &flat)
	flatId := int(flat["id"].(float64))

	// Step 5: Update flat status
	updatePayload := map[string]interface{}{
		"id":     flatId,
		"status": "approved",
	}
	updateBody, _ := json.Marshal(updatePayload)
	req, _ = http.NewRequest("POST", "http://localhost:8080/flat/update", bytes.NewBuffer(updateBody))
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Step 6: Get all flats
	req, _ = http.NewRequest("GET", fmt.Sprintf("http://localhost:8080/house/%v", houseID), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = client.Do(req)
	respHoused, _ := io.ReadAll(resp.Body)
	var housee []map[string]interface{}
	json.Unmarshal(respHoused, &housee)
	fmt.Println(housee)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Define the expected data structure
	expectedData := []map[string]interface{}{
		{"house_id": float64(houseID), "number": float64(1), "price": float64(100000), "rooms": float64(2), "status": "created"},
		{"house_id": float64(houseID), "number": float64(2), "price": float64(200000), "rooms": float64(2), "status": "created"},
		{"house_id": float64(houseID), "number": float64(3), "price": float64(300000), "rooms": float64(2), "status": "approved"},
	}

	// Check for expected values
	for i, flat := range housee {
		assert.Equal(t, expectedData[i]["house_id"], flat["house_id"])
		assert.Equal(t, expectedData[i]["number"], flat["number"])
		assert.Equal(t, expectedData[i]["price"], flat["price"])
		assert.Equal(t, expectedData[i]["rooms"], flat["rooms"])
		assert.Equal(t, expectedData[i]["status"], flat["status"])
	}
}
