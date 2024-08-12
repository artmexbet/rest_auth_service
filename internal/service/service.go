package service

import (
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"net/http"
	"restAuthPart/internal/models"
)

type IJWTManager interface {
	GenerateToken(guid uuid.UUID, ip string) (string, error)
}

type IDatabase interface {
	AddUser(user models.User) error
	AddRefreshToken(token string, guid uuid.UUID) error
}

// Service ...
type Service struct {
	jwtManager IJWTManager
	db         IDatabase
}

// New ...
func New(manager IJWTManager, db IDatabase) *Service {
	return &Service{
		jwtManager: manager,
		db:         db,
	}
}

// Auth returns http.HandlerFunc which process the auth request
// gets from request GUID and returns Access and Refresh tokens
// to ResponseWriter
func (s *Service) Auth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		guidString := chi.URLParam(r, "guid")
		guid, err := uuid.Parse(guidString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// TODO: implement
	}
}

// Refresh returns http.HandlerFunc which process the refresh request
// gets Refresh token and returns new Access and Refresh tokens
func (s *Service) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement
	}
}
