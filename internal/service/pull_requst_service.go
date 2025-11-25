package service

import (
	"context"
	"time"

	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/models"
)

type PullRequestRepository interface {
	CreatePRWithReviewers(ctx context.Context, pr *models.PullRequest, reviewerIDs []string) error
	GetPRByID(ctx context.Context, prID string) (*models.PullRequest, error)
	MergePR(ctx context.Context, prID string) (*models.PullRequest, error)
	ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error
	GetUserReviews(ctx context.Context, userID string) ([]models.PullRequestShort, error)
}

type PullRequestService struct {
	prRepo   PullRequestRepository
	userRepo UserRepository
	teamRepo TeamRepository
}

func NewPullRequestService(
	prRepo PullRequestRepository,
	userRepo UserRepository,
	teamRepo TeamRepository,
) *PullRequestService {
	return &PullRequestService{
		prRepo:   prRepo,
		userRepo: userRepo,
		teamRepo: teamRepo,
	}
}

func (s *PullRequestService) CreatePR(
	ctx context.Context,
	req *models.PostPullRequestCreateJSONBody,
) (*models.PullRequest, error) {

	pr := &models.PullRequest{
		PullRequestShort: models.PullRequestShort{
			PullRequestId:   req.PullRequestId,
			PullRequestName: req.PullRequestName,
			AuthorId:        req.AuthorId,
			Status:          models.PullRequestStatusOPEN,
		},
		CreatedAt: ptr(time.Now()),
	}

	teamID, isActive, err := s.userRepo.GetUserTeamIDAndActive(ctx, req.AuthorId)
	if err != nil {
		return nil, err
	}
	if !isActive {
		return nil, err // TODO: models.ErrNoCandidate
	}

	excludeUserIDs := []string{req.AuthorId}
	reviewerIDs, err := s.userRepo.GetActiveUsersForReview(ctx, teamID, excludeUserIDs)
	if err != nil {
		return nil, err
	}

	if err := s.prRepo.CreatePRWithReviewers(ctx, pr, reviewerIDs); err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewerIDs
	return pr, nil
}

func (s *PullRequestService) ReassignReviewer(
	ctx context.Context,
	prID, oldReviewerID string,
) (*models.ReassignResponse, error) {

	pr, err := s.prRepo.GetPRByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == models.PullRequestStatusMERGED {
		return nil, models.NewAppError(models.PR_MERGED)
	}

	found := false
	for _, r := range pr.AssignedReviewers {
		if r == oldReviewerID {
			found = true
			break
		}
	}
	if !found {
		return nil, models.NewAppError(models.NOT_ASSIGNED)
	}

	teamID, isActive, err := s.userRepo.GetUserTeamIDAndActive(ctx, oldReviewerID)
	if err != nil || !isActive {
		return nil, err
	}

	var secondReviewer string
	for _, r := range pr.AssignedReviewers {
		if r != oldReviewerID {
			secondReviewer = r
			break
		}
	}

	excludeUserIDs := []string{pr.AuthorId, oldReviewerID}
	if secondReviewer != "" {
		excludeUserIDs = append(excludeUserIDs, secondReviewer)
	}

	candidates, err := s.userRepo.GetActiveUsersForReview(ctx, teamID, excludeUserIDs)
	if err != nil {
		return nil, err
	}
	if len(candidates) == 0 {
		return nil, models.NewAppError(models.NO_CANDIDATE)
	}

	newReviewerID := candidates[0]

	if err := s.prRepo.ReplaceReviewer(ctx, prID, oldReviewerID, newReviewerID); err != nil {
		return nil, err
	}

	for i, r := range pr.AssignedReviewers {
		if r == oldReviewerID {
			pr.AssignedReviewers[i] = newReviewerID
			break
		}
	}

	return &models.ReassignResponse{
		PR:         pr,
		ReplacedBy: newReviewerID,
	}, nil
}

func (s *PullRequestService) MergePR(ctx context.Context, prID string) (*models.PullRequest, error) {
	pr, err := s.prRepo.MergePR(ctx, prID)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func (s *PullRequestService) GetUserReviews(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	return s.prRepo.GetUserReviews(ctx, userID)
}

func ptr[T any](v T) *T {
	return &v
}
