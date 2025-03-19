package config

import (
	"cmp"
	"os"
)

func Load() *Config {
	return &Config{
		HTTP: HTTP{
			Port: cmp.Or(os.Getenv("HTTP_PORT"), ":8080"),
		},
		DB: DB{
			Driver: cmp.Or(os.Getenv("DB_DRIVER"), "sqlite3"),
			DSN:    cmp.Or(os.Getenv("DB_DSN"), "./db.sqlite"),
		},
	}
}

type Config struct {
	HTTP HTTP
	DB   DB
}

type HTTP struct {
	Port string
}

type DB struct {
	Driver string
	DSN    string
}
