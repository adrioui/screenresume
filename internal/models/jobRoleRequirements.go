package models

type JobRoleRequirements struct {
	JobRoleID          string  `json:"job_role_id"`
	SkillID            string  `json:"skill_id"`
	Required           bool    `json:"required"`
	MinExperienceYears float64 `json:"min_experience_years"`
	Importance         int     `json:"importance"`
}

type JobRoleRequirementsCreate struct {
	JobRoleID          string  `json:"job_role_id"`
	SkillID            string  `json:"skill_id"`
	Required           bool    `json:"required"`
	MinExperienceYears float64 `json:"min_experience_years"`
	Importance         int     `json:"importance"`
}

type JobRoleRequirementsUpdate struct {
	Required           *bool    `json:"required,omitempty"`
	MinExperienceYears *float64 `json:"min_experience_years,omitempty"`
	Importance         *int     `json:"importance,omitempty" binding:"omitempty,min=1,max=5"`
}
