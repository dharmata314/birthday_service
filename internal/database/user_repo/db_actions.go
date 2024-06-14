package database

import (
	"birthday-service/internal/entities"
	errMsg "birthday-service/internal/err"
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewUserRepository(db *pgxpool.Pool, log *slog.Logger) *UserRepository {
	return &UserRepository{db: db, log: log}
}

func (u *UserRepository) CreateUser(ctx context.Context, user *entities.User) error {
	err := u.db.QueryRow(ctx, `INSERT INTO Users (email, password) VALUES ($1, $2) RETURNING id`, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		u.log.Error("Failed to create user", errMsg.Err(err))
		return err
	}
	return nil
}

func (u *UserRepository) FindUserByEmail(ctx context.Context, email string) (entities.User, error) {
	query, err := u.db.Query(ctx, `SELECT id, email, password FROM Users WHERE email = $1`, email)
	if err != nil {
		u.log.Error("Error querying users table", errMsg.Err(err))
		return entities.User{}, err
	}
	row := entities.User{}
	defer query.Close()
	if !query.Next() {
		u.log.Error("user not found")
		return entities.User{}, fmt.Errorf("user not found")
	} else {
		err := query.Scan(&row.ID, &row.Email, &row.Password)
		if err != nil {
			u.log.Error("Error scanning users", errMsg.Err(err))
			return entities.User{}, err
		}
	}
	return row, nil
}

func (u *UserRepository) FindUserById(ctx context.Context, id int) (entities.User, error) {
	query, err := u.db.Query(ctx, `SELECT id, email, password FROM Users WHERE id = $1`, id)
	if err != nil {
		u.log.Error("error querying users", errMsg.Err(err))
		return entities.User{}, err
	}
	defer query.Close()
	rowArray := entities.User{}
	if !query.Next() {
		u.log.Error("user not found")
		return entities.User{}, fmt.Errorf("user not found")
	} else {
		err := query.Scan(&rowArray.ID, &rowArray.Email, &rowArray.Password)
		if err != nil {
			u.log.Error("error scanning users", errMsg.Err(err))
			return entities.User{}, err
		}
	}
	return rowArray, nil
}

func (u *UserRepository) DeleteUserById(ctx context.Context, id int) error {
	_, err := u.db.Exec(ctx, `DELETE FROM Users WHERE id = $1`, id)
	if err != nil {
		u.log.Error("failed to delete user", errMsg.Err(err))
		return err
	}
	return nil
}

func (u *UserRepository) UpdateUser(ctx context.Context, user *entities.User) error {
	_, err := u.db.Exec(ctx, `UPDATE Users SET email = $1, password = $2 WHERE id = $3`, user.Email, user.Password, user.ID)
	if err != nil {
		u.log.Error("failed to update user", errMsg.Err(err))
		return err
	}

	return nil

}
