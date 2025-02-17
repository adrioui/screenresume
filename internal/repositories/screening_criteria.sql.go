// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: screening_criteria.sql

package repositories

import (
	"context"

	"github.com/google/uuid"
)

const createScreeningCriteria = `-- name: CreateScreeningCriteria :one
insert into screening_criteria (
  id, screening_result_id, decision, reasoning, matched_skills, missing_skills
) values (
  uuid_generate_v4(), $1, $2, $3, $4, $5
)
RETURNING id, screening_result_id, decision, reasoning, matched_skills, missing_skills
`

type CreateScreeningCriteriaParams struct {
	ScreeningResultID uuid.UUID   `json:"screening_result_id"`
	Decision          bool        `json:"decision"`
	Reasoning         string      `json:"reasoning"`
	MatchedSkills     []uuid.UUID `json:"matched_skills"`
	MissingSkills     []uuid.UUID `json:"missing_skills"`
}

func (q *Queries) CreateScreeningCriteria(ctx context.Context, arg CreateScreeningCriteriaParams) (ScreeningCriterium, error) {
	row := q.db.QueryRow(ctx, createScreeningCriteria,
		arg.ScreeningResultID,
		arg.Decision,
		arg.Reasoning,
		arg.MatchedSkills,
		arg.MissingSkills,
	)
	var i ScreeningCriterium
	err := row.Scan(
		&i.ID,
		&i.ScreeningResultID,
		&i.Decision,
		&i.Reasoning,
		&i.MatchedSkills,
		&i.MissingSkills,
	)
	return i, err
}

const getScreeningCriteria = `-- name: GetScreeningCriteria :one
select id, screening_result_id, decision, reasoning, matched_skills, missing_skills
from screening_criteria
where id = $1
limit 1
`

func (q *Queries) GetScreeningCriteria(ctx context.Context, id uuid.UUID) (ScreeningCriterium, error) {
	row := q.db.QueryRow(ctx, getScreeningCriteria, id)
	var i ScreeningCriterium
	err := row.Scan(
		&i.ID,
		&i.ScreeningResultID,
		&i.Decision,
		&i.Reasoning,
		&i.MatchedSkills,
		&i.MissingSkills,
	)
	return i, err
}

const listScreeningCriteria = `-- name: ListScreeningCriteria :many
select id, screening_result_id, decision, reasoning, matched_skills, missing_skills
from screening_criteria
`

func (q *Queries) ListScreeningCriteria(ctx context.Context) ([]ScreeningCriterium, error) {
	rows, err := q.db.Query(ctx, listScreeningCriteria)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []ScreeningCriterium{}
	for rows.Next() {
		var i ScreeningCriterium
		if err := rows.Scan(
			&i.ID,
			&i.ScreeningResultID,
			&i.Decision,
			&i.Reasoning,
			&i.MatchedSkills,
			&i.MissingSkills,
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
