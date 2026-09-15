package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port         string
	DBPath       string
	JWTSecret    string
	UploadDir    string
	Environment  string
	AllowedHosts []string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "opsdesk.db"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "opsdesk_internal_production_signing_key_2026_secured"
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./uploads"
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	return &Config{
		Port:         port,
		DBPath:       dbPath,
		JWTSecret:    jwtSecret,
		UploadDir:    uploadDir,
		Environment:  env,
		AllowedHosts: []string{"http://localhost:3000", "http://127.0.0.1:3000"},
	}
}

func GetEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultVal
	}
	return val
}
