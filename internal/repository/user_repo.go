package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/models"
)

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	var u models.User
	err := r.db.QueryRowContext(ctx,
		`UPDATE users
		 SET is_active = $1
		 WHERE id = $2
		 RETURNING id, username,
		           (SELECT name FROM teams WHERE id = users.team_id) AS team_name,
		           is_active`,
		isActive, userID,
	).Scan(&u.UserId, &u.Username, &u.TeamName, &u.IsActive)

	if err == sql.ErrNoRows {
		return nil, models.NewAppError(models.NOT_FOUND)
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepo) GetUserTeamIDAndActive(ctx context.Context, userID string) (teamID int, isActive bool, err error) {
	err = r.db.QueryRowContext(ctx,
		`SELECT team_id, is_active FROM users WHERE id = $1`,
		userID,
	).Scan(&teamID, &isActive)
	if err == sql.ErrNoRows {
		err = models.NewAppError(models.NOT_FOUND)
	}
	return
}

func (r *userRepo) GetActiveUsersForReview(ctx context.Context, teamID int, excludeUserIDs []string) ([]string, error) {
	query := `
        SELECT id 
        FROM users 
        WHERE team_id = $1 
          AND is_active = true`

	args := []any{teamID}
	if len(excludeUserIDs) > 0 {
		placeholders := make([]string, len(excludeUserIDs))
		for i := range excludeUserIDs {
			placeholders[i] = fmt.Sprintf("$%d", len(args)+1)
			args = append(args, excludeUserIDs[i])
		}
		query += fmt.Sprintf(" AND id NOT IN (%s)", strings.Join(placeholders, ","))
	}

	query += ` ORDER BY RANDOM() LIMIT 2`

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var userIDs []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, id)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}
