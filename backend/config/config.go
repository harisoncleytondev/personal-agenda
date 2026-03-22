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

func GetSMTPHost() string {
	return os.Getenv("SMTP_HOST")
}

func GetSMTPPort() string {
	return os.Getenv("SMTP_PORT")
}

func GetSMTPUser() string {
	return os.Getenv("SMTP_USER")
}

func GetSMTPPass() string {
	return os.Getenv("SMTP_PASS")
}

func GetSMTPFrom() string {
	return os.Getenv("SMTP_FROM")
}