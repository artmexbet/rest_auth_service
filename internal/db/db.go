package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	_ "golang.org/x/crypto/bcrypt"
	"restAuthPart/internal/models"
)

// Config ...
type Config struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	User     string `yaml:"user" env:"USER" env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
	DbName   string `yaml:"dbName" env:"DB_NAME" env-default:"testtask"`
}

// DB ...
type DB struct {
	cfg *Config
	db  *pgx.Conn
}

// New ...
func New(cfg *Config) (*DB, error) {
	d := &DB{
		cfg: cfg,
	}

	db, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DbName))
	if err != nil {
		return nil, err
	}
	d.db = db

	return d, nil
}

// Close ...
func (d *DB) Close() error {
	return d.db.Close(context.Background())
}

// AddUserIfNotExist insert models.User to table users and update its ip if it exists
func (d *DB) AddUserIfNotExist(user models.User) error {
	_, err := d.db.Exec(context.Background(),
		`INSERT INTO public.users (id, ip) 
			 VALUES ($1, $2)
			 ON CONFLICT (id) DO UPDATE SET ip=$2`, user.Guid, user.Ip,
	)
	return err
}

// AddRefreshToken insert token for user to table tokens
func (d *DB) AddRefreshToken(token string, guid uuid.UUID) (int, error) {
	var id int

	//bytes, err := bcrypt.GenerateFromPassword([]byte(token), 8) // Use cost 8 because it is not stated in task
	//if err != nil {
	//	return 0, err
	//}
	err := d.db.QueryRow(context.Background(),
		`INSERT INTO public.tokens (user_id, token)
			 VALUES ($1, $2) RETURNING id`, guid, []byte(token)).Scan(&id)
	return id, err
}

// GetRefreshToken returns hashed token from table tokens
func (d *DB) GetRefreshToken(refreshTokenId int) ([]byte, error) {
	var hashedToken []byte

	err := d.db.QueryRow(context.Background(),
		`SELECT token FROM public.tokens WHERE id=$1`, refreshTokenId).Scan(&hashedToken)
	return hashedToken, err
}
