package repository

import (
	"context"
	"database/sql"

	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/models"
)

type teamRepo struct {
	db *sql.DB
}

func NewTeamRepo(db *sql.DB) *teamRepo {
	return &teamRepo{db: db}
}

func (r *teamRepo) UpsertTeam(ctx context.Context, team *models.Team) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var teamID int
	err = tx.QueryRowContext(ctx,
		`INSERT INTO teams (name) VALUES ($1) RETURNING id`,
		team.TeamName,
	).Scan(&teamID)
	if err != nil {
		return err
	}

	for _, m := range team.Members {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO users (user_id, username, team_id, is_active)
			 VALUES ($1, $2, $3, $4)`,
			m.UserId, m.Username, teamID, m.IsActive,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *teamRepo) GetTeamByName(ctx context.Context, teamName string) (*models.Team, error) {
	team := &models.Team{
		TeamName: teamName,
		Members:  []models.TeamMember{},
	}

	var teamID int
	err := r.db.QueryRowContext(ctx, `SELECT id FROM teams WHERE name = $1`, teamName).Scan(&teamID)
	if err == sql.ErrNoRows {
		return nil, err // TODO: models.NOT_FOUND
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT user_id, username, is_active FROM users WHERE team_id = $1 ORDER BY username`,
		teamID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var m models.TeamMember
		if err := rows.Scan(&m.UserId, &m.Username, &m.IsActive); err != nil {
			return nil, err
		}
		team.Members = append(team.Members, m)
	}

	return team, rows.Err()
}

func (r *teamRepo) GetTeamIDByName(ctx context.Context, teamName string) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx, `SELECT id FROM teams WHERE name = $1`, teamName).Scan(&id)
	if err == sql.ErrNoRows {
		return 0, err //TODO: models.NOT_FOUND
	}
	if err != nil {
		return 0, err
	}
	return id, nil
}
