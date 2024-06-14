package database

import (
	"birthday-service/internal/entities"
	errMsg "birthday-service/internal/err"
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SubsRepository struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewSubsRepository(db *pgxpool.Pool, log *slog.Logger) *SubsRepository {
	return &SubsRepository{db, log}
}

func (s *SubsRepository) CreateSub(ctx context.Context, sub *entities.Subscription) error {
	err := s.db.QueryRow(ctx, `INSERT INTO Subscriptions (user_id, emp_id) VALUES ($1, $2) RETURNING ID`, sub.UserID, sub.EmployeeID).Scan(&sub.ID)
	if err != nil {
		s.log.Error("failed to create subscription", errMsg.Err(err))
		return err
	}
	return nil
}

func (s *SubsRepository) DeleteSub(ctx context.Context, id int) error {
	_, err := s.db.Exec(ctx, `DELETE FROM Subscriptions WHERE id = $1`, id)
	if err != nil {
		s.log.Error("failed to delete user", errMsg.Err(err))
		return err
	}
	return nil
}

func (s *SubsRepository) GetSubs(ctx context.Context, EmployeeID int) ([]entities.User, error) {
	var users []entities.User
	query := `SELECT u.id, u.email
		FROM Users u
		JOIN Subscriptions s ON u.id = s.user_id
		WHERE s.emp_id = $1`

	rows, err := s.db.Query(ctx, query, EmployeeID)
	if err != nil {
		s.log.Error("failed to get subscribers", errMsg.Err(err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user entities.User
		if err := rows.Scan(&user.ID, &user.Email); err != nil {
			s.log.Error("failed to scan user", errMsg.Err(err))
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
