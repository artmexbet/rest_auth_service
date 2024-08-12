package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"restAuthPart/internal/models"
)

type Config struct {
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	User     string `yaml:"user" env:"USER" env-default:"postgres"`
	Password string `yaml:"password" env:"PASSWORD" env-default:"postgres"`
	DbName   string `yaml:"dbName" env:"DB_NAME" env-default:"testtask"`
}

type DB struct {
	cfg *Config
	db  *pgx.Conn
}

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

func (d *DB) Close() error {
	return d.db.Close(context.Background())
}

func (d *DB) AddUser(user models.User) error {
	_, err := d.db.Exec(context.Background(),
		`INSERT INTO public.users (id, ip) 
			 VALUES ($1, $2)`, user.Guid, user.Ip,
	)
	return err
}

func (d *DB) AddRefreshToken(token string, guid uuid.UUID) error {
	_, err := d.db.Exec(context.Background(),
		`INSERT INTO public.tokens (user_id, token)
			 VALUES ($1, $2)`, guid, token)
	return err
}
