package models

import "time"

type Candidates struct {
	ID       string `json:"id"`
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FileID   string `json:"file_id"`
	Status   string `json:"status"`
}

type CandidatesCreate struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FileID   string `json:"file_id"`
	Status   string `json:"status"`
}

type CandidatesUpdate struct {
	FullName string `json:"full_name"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FileID   string `json:"file_id"`
	Status   string `json:"status"`
}

type CandidateAndJobRoles struct {
	Candidate Candidates `json:"candidate"`
	JobRole   JobRoles   `json:"job_role"`
	AppliedAt time.Time  `json:"applied_at"`
}
