package repository

import (
	"context"
	"cruder/internal/model"
	"database/sql"
)

type UserRepository interface {
	GetAll() ([]model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByID(id int64) (*model.User, error)
	GetByUUID(uuid string) (*model.User, error)
	Create(user *model.User) error
	Update(uuid string, user *model.User) error
	Delete(uuid string) error
}

type userRepository struct {
	db *sql.DB
}

const selectUserColumns = "SELECT id, uuid, username, email, full_name FROM users"

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (userRepository *userRepository) buildSelectQuery(whereClause string) string {
	return selectUserColumns + " WHERE " + whereClause
}

func (userRepository *userRepository) scanUserRow(row *sql.Row) (*model.User, error) {
	var user model.User
	err := row.Scan(&user.ID, &user.UUID, &user.Username, &user.Email, &user.FullName)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (userRepository *userRepository) GetAll() ([]model.User, error) {
	rows, err := userRepository.db.QueryContext(context.Background(), selectUserColumns)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.UUID, &user.Username, &user.Email, &user.FullName); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, rows.Err()
}

func (userRepository *userRepository) GetByUsername(username string) (*model.User, error) {
	row := userRepository.db.QueryRowContext(context.Background(), userRepository.buildSelectQuery("username = $1"), username)
	return userRepository.scanUserRow(row)
}

func (userRepository *userRepository) GetByID(id int64) (*model.User, error) {
	row := userRepository.db.QueryRowContext(context.Background(), userRepository.buildSelectQuery("id = $1"), id)
	return userRepository.scanUserRow(row)
}

func (userRepository *userRepository) GetByUUID(uuid string) (*model.User, error) {
	row := userRepository.db.QueryRowContext(context.Background(), userRepository.buildSelectQuery("uuid = $1"), uuid)
	return userRepository.scanUserRow(row)
}

func (userRepository *userRepository) Create(user *model.User) error {
	query := `INSERT INTO users (username, email, full_name) VALUES ($1, $2, $3) RETURNING id, uuid`
	err := userRepository.db.QueryRowContext(context.Background(), query, user.Username, user.Email, user.FullName).
		Scan(&user.ID, &user.UUID)
	return err
}

func (userRepository *userRepository) Update(uuid string, user *model.User) error {
	query := `UPDATE users SET username = $1, email = $2, full_name = $3 WHERE uuid = $4`
	_, err := userRepository.db.ExecContext(context.Background(), query, user.Username, user.Email, user.FullName, uuid)
	return err
}

func (userRepository *userRepository) Delete(uuid string) error {
	query := `DELETE FROM users WHERE uuid = $1`
	_, err := userRepository.db.ExecContext(context.Background(), query, uuid)
	return err
}
