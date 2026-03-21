package model

import "time"

type User struct {
    ID           string    `json:"id" db:"id"`
    Name         string    `json:"name" db:"name"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"password_hash" db:"password_hash"`
    CreatedAt    time.Time `json:"created_at" db:"created_at"`
}

type UserCreate struct {
    Name         string    `json:"name" db:"name"`
    Email        string    `json:"email" db:"email"`
    PasswordHash string    `json:"password_hash" db:"password_hash"`
}