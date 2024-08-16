package service

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"restAuthPart/internal/models"
)

type IJWTManager interface {
	GenerateRefreshToken(guid uuid.UUID, ip string) (string, error)
	GenerateAccessToken(guid uuid.UUID, ip string, id int) (string, error)
	GetClaims(token string, claimsType jwt.Claims) (jwt.Claims, error)
	CompareTokens(token string, hashedToken []byte) bool
}

type IDatabase interface {
	AddUserIfNotExist(user models.User) error
	AddRefreshToken(token string, guid uuid.UUID) (int, error)
	GetRefreshToken(refreshTokenId int) ([]byte, error)
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

		if err := s.db.AddUserIfNotExist(models.User{Guid: guid, Ip: r.RemoteAddr}); err != nil {
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
		logger := slog.With(slog.String("module", "Service.Refresh"))

		var data models.RefreshTokenJSON
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, "Can't parse json", http.StatusBadRequest)
			logger.Error("Can't parse json", slog.String("err", err.Error()))
			return
		}

		refreshClaims, err := s.jwtManager.GetClaims(data.RefreshT, &models.RefreshTokenClaims{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("Cannot get refreshClaims from token", slog.String("err", err.Error()))
			return
		}

		decodedRefreshClaims, ok := refreshClaims.(*models.RefreshTokenClaims)
		if !ok {
			http.Error(w, "Cannot convert refreshClaims to RefreshTokenClaims", http.StatusBadRequest)
			logger.Error("Cannot convert refreshClaims to RefreshTokenClaims")
			return
		}

		accessClaims, err := s.jwtManager.GetClaims(data.AccessT, &models.AccessTokenClaims{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("Cannot get accessClaims from token", slog.String("err", err.Error()))
			return
		}

		fmt.Println(accessClaims)
		decodedAccessClaims, ok := accessClaims.(*models.AccessTokenClaims)

		if !ok {
			http.Error(w, "Cannot convert accessClaims to AccessTokenClaims", http.StatusBadRequest)
			logger.Error("Cannot convert accessClaims to AccessTokenClaims")
			return
		}

		if decodedRefreshClaims.Guid != decodedAccessClaims.Guid {
			http.Error(w, "Guid from token doesn't match guid from request", http.StatusBadRequest)
			logger.Error("Guid from token doesn't match guid from request")
			return
		}

		tokenFromDb, err := s.db.GetRefreshToken(decodedAccessClaims.RefreshId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			logger.Error("Cannot get token from DB", slog.String("err", err.Error()))
			return
		}

		if !s.jwtManager.CompareTokens(data.RefreshT, tokenFromDb) {
			http.Error(w, "Tokens are not identical", http.StatusBadRequest)
			logger.Error("Tokens are not identical")
			return
		}

		if decodedRefreshClaims.Ip != r.RemoteAddr {
			// TODO: Handle this case
			//  Send email to address
		}

		accessToken, err := s.jwtManager.GenerateAccessToken(decodedRefreshClaims.Guid, r.RemoteAddr, decodedAccessClaims.RefreshId)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Cannot generate access token", slog.String("err", err.Error()))
			return
		}

		w.WriteHeader(http.StatusAccepted)
		if err := json.NewEncoder(w).Encode(models.AccessRefreshJSON{AccessT: accessToken, RefreshT: data.RefreshT}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Error("Cannot write encoded json", slog.String("err", err.Error()))
		}
	}
}
