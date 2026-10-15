package config

import (
	"os"
	"time"
)

type Config struct {
	Addr      string
	DBPath    string
	JWTSecret string
	TokenTTL  time.Duration
	StaticDir string
}

func Load() Config {
	return Config{
		Addr:      env("APP_ADDR", ":8080"),
		DBPath:    env("APP_DB", "data/admin.db"),
		JWTSecret: env("APP_JWT_SECRET", "change-me-in-production"),
		TokenTTL:  24 * time.Hour,
		StaticDir: env("APP_STATIC_DIR", "../web/dist"),
	}
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
