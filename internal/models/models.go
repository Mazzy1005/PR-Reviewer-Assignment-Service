package models

import (
	"time"
)

type ErrorResponseErrorCode string

const (
	TEAMEXISTS  ErrorResponseErrorCode = "TEAM_EXISTS"
	PREXISTS    ErrorResponseErrorCode = "PR_EXISTS"
	PRMERGED    ErrorResponseErrorCode = "PR_MERGED"
	NOTASSIGNED ErrorResponseErrorCode = "NOT_ASSIGNED"
	NOCANDIDATE ErrorResponseErrorCode = "NO_CANDIDATE"
	NOTFOUND    ErrorResponseErrorCode = "NOT_FOUND"
)

type ErrorResponse struct {
	Error struct {
		Code    ErrorResponseErrorCode `json:"code"`
		Message string                 `json:"message"`
	} `json:"error"`
}

type PullRequestStatus string

const (
	PullRequestStatusMERGED PullRequestStatus = "MERGED"
	PullRequestStatusOPEN   PullRequestStatus = "OPEN"
)

type PullRequestShort struct {
	PullRequestId   string            `json:"pull_request_id"`
	PullRequestName string            `json:"pull_request_name"`
	AuthorId        string            `json:"author_id"`
	Status          PullRequestStatus `json:"status"`
}

type PullRequest struct {
	PullRequestShort
	AssignedReviewers []string   `json:"assigned_reviewers"`
	CreatedAt         *time.Time `json:"createdAt"`
	MergedAt          *time.Time `json:"mergedAt"`
}

type TeamMember struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	IsActive bool   `json:"is_active"`
}

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

type User struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
	IsActive bool   `json:"is_active"`
}

type PostPullRequestCreateJSONBody struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
}

type PostPullRequestMergeJSONBody struct {
	PullRequestId string `json:"pull_request_id"`
}

type PostPullRequestReassignJSONBody struct {
	PullRequestId string `json:"pull_request_id"`
	OldUserId     string `json:"old_user_id"`
}

type GetTeamGetParams struct {
	TeamName string `form:"team_name" json:"team_name"`
}

type GetUsersGetReviewParams struct {
	UserId string `form:"user_id" json:"user_id"`
}

type PostUsersSetIsActiveJSONBody struct {
	UserId   string `json:"user_id"`
	IsActive bool   `json:"is_active"`
}
