package service

import (
	"context"
	"errors"

	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/models"
)

type TeamRepository interface {
	UpsertTeam(ctx context.Context, team *models.Team) error
	GetTeamByName(ctx context.Context, teamName string) (*models.Team, error)
	GetTeamIDByName(ctx context.Context, teamName string) (int, error)
}

type TeamService struct {
	teamRepo TeamRepository
}

func NewTeamService(teamRepo TeamRepository) *TeamService {
	return &TeamService{teamRepo: teamRepo}
}

func (s *TeamService) AddTeam(ctx context.Context, team *models.Team) error {
	if team.TeamName == "" {
		return errors.New("Empty Name") // TODO добавить ошибку
	}
	// Репозиторий сам сделает upsert + транзакцию
	return s.teamRepo.UpsertTeam(ctx, team)
}

// GetTeam — /team/get
func (s *TeamService) GetTeam(ctx context.Context, teamName string) (*models.Team, error) {
	if teamName == "" {
		return nil, errors.New("Empty Name") // TODO добавить ошибку
	}
	return s.teamRepo.GetTeamByName(ctx, teamName)
}
