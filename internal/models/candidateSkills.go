package models

type CandidateSkills struct {
	CandidateID     string  `json:"candidate_id"`
	SkillID         string  `json:"skill_id"`
	YearsExperience float64 `json:"years_experience"`
	LastUsed        string  `json:"last_used"`
}

type CandidateSkillsCreate struct {
	CandidateID     string  `json:"candidate_id"`
	SkillID         string  `json:"skill_id"`
	YearsExperience float64 `json:"years_experience"`
	LastUsed        string  `json:"last_used"`
}

type CandidateSkillsUpdate struct {
	YearsExperience float64 `json:"years_experience"`
	LastUsed        string  `json:"last_used"`
}
