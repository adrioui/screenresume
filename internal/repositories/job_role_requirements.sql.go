// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: job_role_requirements.sql

package repositories

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

const createJobRoleRequirement = `-- name: CreateJobRoleRequirement :one
INSERT INTO job_role_requirements (
  job_role_id, skill_id, required, min_experience_years, importance
) VALUES (
  $1, $2, $3, $4, $5
)
RETURNING job_role_id, skill_id, required, min_experience_years, importance
`

type CreateJobRoleRequirementParams struct {
	JobRoleID          uuid.UUID      `json:"job_role_id"`
	SkillID            uuid.UUID      `json:"skill_id"`
	Required           bool           `json:"required"`
	MinExperienceYears pgtype.Numeric `json:"min_experience_years"`
	Importance         int32          `json:"importance"`
}

func (q *Queries) CreateJobRoleRequirement(ctx context.Context, arg CreateJobRoleRequirementParams) (JobRoleRequirement, error) {
	row := q.db.QueryRow(ctx, createJobRoleRequirement,
		arg.JobRoleID,
		arg.SkillID,
		arg.Required,
		arg.MinExperienceYears,
		arg.Importance,
	)
	var i JobRoleRequirement
	err := row.Scan(
		&i.JobRoleID,
		&i.SkillID,
		&i.Required,
		&i.MinExperienceYears,
		&i.Importance,
	)
	return i, err
}

const deleteJobRoleRequirement = `-- name: DeleteJobRoleRequirement :exec
delete from job_role_requirements
where job_role_id = $1 and skill_id = $2
`

type DeleteJobRoleRequirementParams struct {
	JobRoleID uuid.UUID `json:"job_role_id"`
	SkillID   uuid.UUID `json:"skill_id"`
}

func (q *Queries) DeleteJobRoleRequirement(ctx context.Context, arg DeleteJobRoleRequirementParams) error {
	_, err := q.db.Exec(ctx, deleteJobRoleRequirement, arg.JobRoleID, arg.SkillID)
	return err
}

const getJobRoleRequirement = `-- name: GetJobRoleRequirement :one
select job_role_id, skill_id, required, min_experience_years, importance
from job_role_requirements
where job_role_id = $1 and skill_id = $2
limit 1
`

type GetJobRoleRequirementParams struct {
	JobRoleID uuid.UUID `json:"job_role_id"`
	SkillID   uuid.UUID `json:"skill_id"`
}

func (q *Queries) GetJobRoleRequirement(ctx context.Context, arg GetJobRoleRequirementParams) (JobRoleRequirement, error) {
	row := q.db.QueryRow(ctx, getJobRoleRequirement, arg.JobRoleID, arg.SkillID)
	var i JobRoleRequirement
	err := row.Scan(
		&i.JobRoleID,
		&i.SkillID,
		&i.Required,
		&i.MinExperienceYears,
		&i.Importance,
	)
	return i, err
}

const listJobRoleRequirements = `-- name: ListJobRoleRequirements :many
select job_role_id, skill_id, required, min_experience_years, importance
from job_role_requirements
order by job_role_id, skill_id
`

func (q *Queries) ListJobRoleRequirements(ctx context.Context) ([]JobRoleRequirement, error) {
	rows, err := q.db.Query(ctx, listJobRoleRequirements)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []JobRoleRequirement{}
	for rows.Next() {
		var i JobRoleRequirement
		if err := rows.Scan(
			&i.JobRoleID,
			&i.SkillID,
			&i.Required,
			&i.MinExperienceYears,
			&i.Importance,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const toggleJobRoleRequirement = `-- name: ToggleJobRoleRequirement :exec
UPDATE job_role_requirements
SET required = NOT required
WHERE job_role_id = $1 AND skill_id = $2
`

type ToggleJobRoleRequirementParams struct {
	JobRoleID uuid.UUID `json:"job_role_id"`
	SkillID   uuid.UUID `json:"skill_id"`
}

func (q *Queries) ToggleJobRoleRequirement(ctx context.Context, arg ToggleJobRoleRequirementParams) error {
	_, err := q.db.Exec(ctx, toggleJobRoleRequirement, arg.JobRoleID, arg.SkillID)
	return err
}

const updateJobRoleRequirement = `-- name: UpdateJobRoleRequirement :exec
UPDATE job_role_requirements
SET required = $3,
    min_experience_years = $4,
    importance = $5
WHERE job_role_id = $1 AND skill_id = $2
`

type UpdateJobRoleRequirementParams struct {
	JobRoleID          uuid.UUID      `json:"job_role_id"`
	SkillID            uuid.UUID      `json:"skill_id"`
	Required           bool           `json:"required"`
	MinExperienceYears pgtype.Numeric `json:"min_experience_years"`
	Importance         int32          `json:"importance"`
}

func (q *Queries) UpdateJobRoleRequirement(ctx context.Context, arg UpdateJobRoleRequirementParams) error {
	_, err := q.db.Exec(ctx, updateJobRoleRequirement,
		arg.JobRoleID,
		arg.SkillID,
		arg.Required,
		arg.MinExperienceYears,
		arg.Importance,
	)
	return err
}
