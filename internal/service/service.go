package service

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"restAuthPart/internal/models"
)

type IJWTManager interface {
	GenerateRefreshToken(guid uuid.UUID, ip string) (string, error)
	GenerateAccessToken(guid uuid.UUID, ip string, refreshId int) (string, error)
}

type IDatabase interface {
	AddUser(user models.User) error
	AddRefreshToken(token string, guid uuid.UUID) (int, error)
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
		// Get guid from request and generate access and refresh tokens for it then
		logger := slog.With(slog.String("module", "Service.Auth"))
		guidString := chi.URLParam(r, "guid")
		guid, err := uuid.Parse(guidString)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("Cannot parse uuid", slog.String("err", err.Error()))
			return
		}

		if err := s.db.AddUser(models.User{Guid: guid, Ip: r.RemoteAddr}); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("Cannot add user", slog.String("err", err.Error()))
			return
		}

		// Generate refresh token and then get th
		rToken, err := s.jwtManager.GenerateRefreshToken(guid, r.RemoteAddr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Cannot generate refresh token", slog.String("err", err.Error()))
			return
		}

		id, err := s.db.AddRefreshToken(rToken, guid)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Cannot add refresh token to DB", slog.String("err", err.Error()))
			return
		}

		accessToken, err := s.jwtManager.GenerateAccessToken(guid, r.RemoteAddr, id)
		if err != nil {
			logger.Error("Cannot generate access token", slog.String("err", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tokenJson := models.AccessRefreshJSON{AccessT: accessToken, RefreshT: rToken}
		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(tokenJson); err != nil {
			logger.Error("Cannot write encoded json", slog.String("err", err.Error()))
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// Refresh returns http.HandlerFunc which process the refresh request
// gets Refresh token and returns new Access and Refresh tokens
func (s *Service) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO: implement
	}
}
