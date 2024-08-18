package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"restAuthPart/internal/models"
	"testing"
	"time"
)

type MockJWTManager struct {
	mock.Mock
}

func (m *MockJWTManager) GenerateRefreshToken(guid uuid.UUID, ip string) (string, error) {
	args := m.Called(guid, ip)
	return args.String(0), args.Error(1)
}

func (m *MockJWTManager) GenerateAccessToken(guid uuid.UUID, ip string, id int) (string, error) {
	args := m.Called(guid, ip, id)
	return args.String(0), args.Error(1)
}

func (m *MockJWTManager) GetClaims(token string, claimsType jwt.Claims) (jwt.Claims, error) {
	args := m.Called(token, claimsType)
	return args.Get(0).(jwt.Claims), args.Error(1)
}
func (m *MockJWTManager) CompareTokens(token string, hashedToken []byte) bool {
	args := m.Called(token, hashedToken)
	return args.Bool(0)
}

type MockDatabase struct {
	mock.Mock
	users  map[uuid.UUID]models.User
	tokens map[int][]byte
}

func (m *MockDatabase) AddUserIfNotExist(user models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockDatabase) AddRefreshToken(token string, guid uuid.UUID) (int, error) {
	args := m.Called(token, guid)
	return args.Int(0), args.Error(1)
}
func (m *MockDatabase) GetRefreshToken(refreshTokenId int) ([]byte, error) {
	args := m.Called(refreshTokenId)
	return args.Get(0).([]byte), args.Error(1)
}
func (m *MockDatabase) GetUser(guid uuid.UUID) (models.User, error) {
	args := m.Called(guid)
	return args.Get(0).(models.User), args.Error(1)
}

type MockEmailService struct {
	mock.Mock
}

func (m *MockEmailService) SendWarning(email string) error {
	args := m.Called(email)
	return args.Error(0)
}

func TestAuth(t *testing.T) {
	// Arrange
	manager := new(MockJWTManager)
	db := new(MockDatabase)
	emailService := new(MockEmailService)
	service := New(manager, db, emailService)

	manager.On("GenerateRefreshToken", mock.Anything, mock.Anything).Return("refreshToken", nil)
	manager.On("GenerateAccessToken", mock.Anything,
		mock.Anything, mock.Anything).Return("accessToken", nil)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.RefreshTokenClaims{
		Guid: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Ip:   "",
	})
	manager.On("GetClaims", mock.Anything, &models.RefreshTokenClaims{}).Return(token.Claims, nil)

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, models.AccessTokenClaims{
		Guid:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Ip:        "",
		RefreshId: 1,
	})
	manager.On("GetClaims", mock.Anything, &models.AccessTokenClaims{}).Return(token.Claims, nil)

	manager.On("CompareTokens", mock.Anything, mock.Anything).Return(true)

	db.On("AddUserIfNotExist", mock.Anything).Return(nil)
	db.On("AddRefreshToken", mock.Anything, mock.Anything).Return(1, nil)
	db.On("GetRefreshToken", mock.Anything).Return([]byte("refreshToken"), nil)
	db.On("GetUser", mock.Anything).Return(models.User{
		Guid:  uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Ip:    "",
		Email: "example@example.com",
	}, nil)

	emailService.On("SendWarning", mock.Anything).Return(nil)

	r := chi.NewRouter()
	r.Get("/{guid}", service.Auth())

	// Act
	// Test good data
	req, _ := http.NewRequest("GET", "/123e4567-e89b-12d3-a456-426614174000", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert
	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}

	expected := `{"accessT":"accessToken","refreshT":"refreshToken"}` + "\n"

	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

	// Test bad data
	req, _ = http.NewRequest("GET", "/123e4567-e89b-", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}

var RefreshToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJndWlkIjoiYzY0M2Y5YjYtMjIwYS00NmFkLWFjYjEtNTkwMmY2NDA1YjY1IiwiaXAiOiIxNzIuMTguMC4xOjQ2MjA2IiwiZXhwIjoxNzI2NTY2NTk5fQ.9hNf-JiePRV-p3ZPGHxfkf7qslVC6qlD4i4Y20Bs31OSn35AE-UN5LX0k7MtRoMtf_f5DDpw7lXQ_UmYOxftHQ"
var AccessToken = "eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJndWlkIjoiYzY0M2Y5YjYtMjIwYS00NmFkLWFjYjEtNTkwMmY2NDA1YjY1IiwiaXAiOiIxNzIuMTguMC4xOjQ2MjA2IiwicmVmcmVzaElkIjo1LCJleHAiOjE3MjM5NzU0OTl9.KH4qZf1qdq5z8kZfd-82Ii4hAizU4dOMbeAgXKNC4SNCICrtUoZ2KIgIvd04qiidoHTeNa82snPaSfWh3gukSw"

func TestRefresh(t *testing.T) {
	// Arrange
	manager := new(MockJWTManager)
	db := new(MockDatabase)
	emailService := new(MockEmailService)
	service := New(manager, db, emailService)

	manager.On("GenerateRefreshToken",
		mock.Anything,
		mock.Anything).Return(
		RefreshToken,
		nil,
	)
	manager.On("GenerateAccessToken",
		mock.Anything,
		mock.Anything,
		mock.Anything).Return(
		AccessToken,
		nil,
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, models.RefreshTokenClaims{
		Guid: uuid.MustParse("c643f9b6-220a-46ad-acb1-5902f6405b65"),
		Ip:   "123",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(60 * 24 * time.Hour)),
		},
	})
	manager.On("GetClaims", RefreshToken, mock.Anything).Return(token.Claims, nil)

	token = jwt.NewWithClaims(jwt.SigningMethodHS256, models.AccessTokenClaims{
		Guid:      uuid.MustParse("c643f9b6-220a-46ad-acb1-5902f6405b65"),
		Ip:        "123",
		RefreshId: 1,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	})
	manager.On("GetClaims", AccessToken, mock.Anything).Return(token.Claims, nil)
	manager.On("GetClaims", mock.Anything, mock.Anything).Return(token.Claims, fmt.Errorf("Cannot get claims"))

	manager.On("CompareTokens", RefreshToken, []byte(RefreshToken)).Return(true)
	manager.On("CompareTokens", mock.Anything, []byte(RefreshToken)).Return(false)

	db.On("AddUserIfNotExist", mock.Anything).Return(nil)
	db.On("AddRefreshToken", mock.Anything, mock.Anything).Return(5, nil)
	db.On("GetRefreshToken", mock.Anything).Return([]byte(RefreshToken), nil)
	db.On("GetUser", mock.Anything).Return(models.User{
		Guid:  uuid.MustParse("c643f9b6-220a-46ad-acb1-5902f6405b65"),
		Ip:    "",
		Email: "example@example.com",
	}, nil)

	emailService.On("SendWarning", mock.Anything).Return(nil)

	r := chi.NewRouter()
	r.Post("/refresh", service.Refresh())

	// Act
	// Test good data
	mockData := models.RefreshTokenJSON{
		RefreshT: RefreshToken,
		AccessT:  AccessToken,
	}
	data, _ := json.Marshal(mockData)
	req, _ := http.NewRequest("POST", "/refresh", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	// Assert
	// TODO: Fix it
	if status := rr.Code; status != http.StatusAccepted {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusAccepted)
	}

	// Test bad data
	// 1
	req, _ = http.NewRequest("POST", "/123e4567-e89b-", nil)
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}

	// 2
	mockData = models.RefreshTokenJSON{
		RefreshT: "refreshToken",
		AccessT:  "accessToken",
	}
	data, _ = json.Marshal(mockData)
	req, _ = http.NewRequest("POST", "/refresh", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	rr = httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusBadRequest)
	}
}
