package repository

import (
	"context"
	"errors"

	"github.com/harisoncleytondev/personal-agenda/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r * UserRepository) UserCreate(ctx context.Context, user * model.UserCreate) error {
	query := `INSERT INTO users (name, email, password_hash) VALUES ($1, $2, $3)`

	_, err := r.db.Exec(ctx, query, user.Name, user.Email, user.PasswordHash)

	if err != nil {
        return err
    }

    return nil
}

func (r * UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, email, password_hash FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).
        Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
    
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, errors.New("usuário não encontrado")
        }
        return nil, err
    }

    return &user, nil
}

func (r * UserRepository) FindById(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	query := `SELECT id, name, email, password_hash FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).
        Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
    
    if err != nil {
        if errors.Is(err, pgx.ErrNoRows) {
            return nil, errors.New("usuário não encontrado")
        }
        return nil, err
    }

    return &user, nil
}