package service

import (
	"context"

	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/models"
)

// type UserRepository interface {
// 	SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
// 	GetActiveUsersByTeamID(ctx context.Context, teamID int) ([]string, error)
// 	GetUserTeamIDAndActive(ctx context.Context, userID string) (teamID int, isActive bool, err error)
// }

type UserRepository interface {
	SetUserActive(ctx context.Context, userID string, isActive bool) (*models.User, error)
	GetUserTeamIDAndActive(ctx context.Context, userID string) (teamID int, isActive bool, err error)
	GetActiveUsersForReview(ctx context.Context, teamID int, excludeUserIDs []string) ([]string, error)
}

type UserService struct {
	userRepo UserRepository
	prRepo   PullRequestRepository
}

func NewUserService(userRepo UserRepository, prRepo PullRequestRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		prRepo:   prRepo,
	}
}

// SetIsActive — /users/setIsActive
func (s *UserService) SetIsActive(ctx context.Context, userID string, isActive bool) (*models.User, error) {
	return s.userRepo.SetUserActive(ctx, userID, isActive)
}

// GetReviews — /users/getReview?user_id=...
func (s *UserService) GetReviews(ctx context.Context, userID string) (*models.GetUserReviewsResponse, error) {
	prs, err := s.prRepo.GetUserReviews(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &models.GetUserReviewsResponse{PullRequests: prs}, nil
}
