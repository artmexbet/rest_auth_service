package main

import (
	_jwt "github.com/golang-jwt/jwt/v5"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"restAuthPart/internal/db"
	"restAuthPart/internal/emailService"
	"restAuthPart/internal/jwt"
	"restAuthPart/internal/logger/sl"
	"restAuthPart/internal/router"
	"restAuthPart/internal/service"
)

// Config ...
type Config struct {
	JWTConfig      jwt.Config    `yaml:"jwt" env-prefix:"JWT_"`
	RouterConfig   router.Config `yaml:"router" env-prefix:"ROUTER_"`
	DatabaseConfig db.Config     `yaml:"db" env-prefix:"DB_"`
}

// readConfig ...
func readConfig(filename string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(filename, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func main() {
	sl.SetupLogger("local")

	cfg, err := readConfig("./config.yml")
	if err != nil {
		log.Fatalln(err)
	}

	database, err := db.New(&cfg.DatabaseConfig)
	if err != nil {
		log.Fatalln(err)
	}

	jwtManager := jwt.New(&cfg.JWTConfig, _jwt.SigningMethodHS512)

	email := emailService.New()

	svc := service.New(jwtManager, database, email)

	r := router.New(&cfg.RouterConfig, svc)
	err = r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
