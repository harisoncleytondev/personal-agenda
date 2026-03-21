package config

import (
	"os"
)

func GetPort() string {
	return os.Getenv("PORT")
}

func getDatabaseLink() string {
	return os.Getenv("DATABASE_PG_LINK")
}

func GetJWTSecret() string {
	return os.Getenv("JWT_SECRET")
}