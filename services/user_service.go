package services

import (
	"context"
	"core_two_go/models"
	"database/sql"
	"fmt"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id int) error
}

type userServiceImpl struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) UserService {
	return &userServiceImpl{db: db}
}

func (s *userServiceImpl) CreateUser(ctx context.Context, user *models.User) error {
	query := `INSERT INTO users (id, name) VALUES ($1, $2)`
	_, err := s.db.ExecContext(ctx, query, user.ID, user.Name)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (s *userServiceImpl) GetUser(ctx context.Context, id int) (*models.User, error) {
	query := `SELECT id, name FROM users WHERE id = $1`
	row := s.db.QueryRowContext(ctx, query, id)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to retrieve user: %w", err)
	}

	return &user, nil
}

func (s *userServiceImpl) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET name = $1 WHERE id = $2`
	result, err := s.db.ExecContext(ctx, query, user.Name, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", user.ID)
	}

	return nil
}

func (s *userServiceImpl) DeleteUser(ctx context.Context, id int) error {
	query := `DELETE FROM users WHERE id = $1`
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no user found with id %d", id)
	}

	return nil
}
