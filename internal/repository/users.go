package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/pkg/errors"
)

type DBUser interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type User struct {
	UUID     string
	Email    string
	Password string
	Type     string
}

type UserRepository struct {
	db DBUser
}

func NewUserRepository(db DBUser) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user User) error {
	if _, err := r.db.Exec(ctx, "INSERT INTO users (email, password, type) VALUES ($1, $2, $3)", user.Email, user.Password, user.Type); err != nil {
		return errors.Wrap(err, "create user")
	}

	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, email string) (User, error) {
	var user User
	if err := r.db.QueryRow(ctx, "SELECT * FROM users WHERE email = $1", email).Scan(&user.UUID, &user.Email, &user.Password, &user.Type); err != nil {
		return User{}, errors.Wrap(err, "get user")
	}

	return user, nil
}
