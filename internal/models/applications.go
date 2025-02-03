package models

import (
	"screenresume/internal/repositories"
	"time"
)

type Application struct {
	ID          string                        `json:"id"`
	CandidateID string                        `json:"candidate_id"`
	JobRoleID   string                        `json:"job_role_id"`
	Stage       repositories.ApplicationStage `json:"stage"`
	Score       float64                       `json:"score"`
	AppliedAt   time.Time                     `json:"applied_at"`
	LastUpdated time.Time                     `json:"last_updated"`
}

type ApplicationCreate struct {
	CandidateID string                        `json:"candidate_id"`
	JobRoleID   string                        `json:"job_role_id"`
	Stage       repositories.ApplicationStage `json:"stage"`
	Score       float64                       `json:"score"`
}

type ApplicationUpdate struct {
	CandidateID string                        `json:"candidate_id"`
	JobRoleID   string                        `json:"job_role_id"`
	Stage       repositories.ApplicationStage `json:"stage"`
	Score       float64                       `json:"score"`
}
