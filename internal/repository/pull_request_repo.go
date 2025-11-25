package repository

import (
	"context"
	"database/sql"

	"github.com/Mazzy1005/PR-Reviewer-Assignment-Service/internal/models"
)

type prRepo struct {
	db *sql.DB
}

func NewPullRequestRepo(db *sql.DB) *prRepo {
	return &prRepo{db: db}
}

func (r *prRepo) CreatePRWithReviewers(ctx context.Context, pr *models.PullRequest, reviewerIDs []string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var existed bool
	err = tx.QueryRowContext(ctx,
		`INSERT INTO pull_requests (id, name, author_id, status)
		 VALUES ($1, $2, $3, 'OPEN')
		 ON CONFLICT (id) DO NOTHING RETURNING true`,
		pr.PullRequestId, pr.PullRequestName, pr.AuthorId,
	).Scan(&existed)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.NewAppError(models.PR_EXISTS)
		}
		return err
	}

	for _, reviewerID := range reviewerIDs {
		_, err = tx.ExecContext(ctx,
			`INSERT INTO pr_reviewers (pr_id, reviewer_id) VALUES ($1, $2)`,
			pr.PullRequestId, reviewerID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
func (r *prRepo) GetPRByID(ctx context.Context, prID string) (*models.PullRequest, error) {
	pr := &models.PullRequest{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, author_id, status, created_at, merged_at
		 FROM pull_requests WHERE id = $1`,
		prID,
	).Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.Status, &pr.CreatedAt, &pr.MergedAt)
	if err == sql.ErrNoRows {
		return nil, models.NewAppError(models.NOT_FOUND)
	}
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT reviewer_id FROM pr_reviewers WHERE pr_id = $1`,
		prID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID string
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		pr.AssignedReviewers = append(pr.AssignedReviewers, userID)
	}
	return pr, rows.Err()
}

func (r *prRepo) MergePR(ctx context.Context, prID string) (*models.PullRequest, error) {
	pr, err := r.GetPRByID(ctx, prID)
	if err != nil {
		return nil, err
	}
	if pr.Status == models.PullRequestStatusMERGED {
		return pr, nil
	}

	_, err = r.db.ExecContext(ctx,
		`UPDATE pull_requests SET status = 'MERGED', merged_at = NOW() WHERE id = $1 AND status = 'OPEN'`,
		prID,
	)
	if err != nil {
		return nil, err
	}

	return r.GetPRByID(ctx, prID)
}

func (r *prRepo) ReplaceReviewer(ctx context.Context, prID, oldReviewerID, newReviewerID string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE pr_reviewers SET reviewer_id = $1 WHERE pr_id = $2 AND reviewer_id = $3`,
		newReviewerID, prID, oldReviewerID,
	)
	return err
}

func (r *prRepo) GetUserReviews(ctx context.Context, userID string) ([]models.PullRequestShort, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT p.id, p.name, p.author_id, p.status
		 FROM pull_requests p
		 JOIN pr_reviewers r ON p.id = r.pr_id
		 WHERE r.reviewer_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prs []models.PullRequestShort
	for rows.Next() {
		var pr models.PullRequestShort
		if err := rows.Scan(&pr.PullRequestId, &pr.PullRequestName, &pr.AuthorId, &pr.Status); err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, rows.Err()
}
