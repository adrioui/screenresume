package models

type ScreeningCriteria struct {
	ID                 string   `json:"id"`
	ScreeningResultsID string   `json:"screening_results_id"`
	CriteriaText       string   `json:"criteria_text"`
	Decision           bool     `json:"decision"`
	Reasoning          string   `json:"reasoning"`
	MatchedSkills      []string `json:"matched_skills"`
	MissingSkills      []string `json:"missing_skills"`
}

type ScreeningCriteriaCreate struct {
	ScreeningResultsID string   `json:"screening_results_id"`
	CriteriaText       string   `json:"criteria_text"`
	Decision           bool     `json:"decision"`
	Reasoning          string   `json:"reasoning"`
	MatchedSkills      []string `json:"matched_skills"`
	MissingSkills      []string `json:"missing_skills"`
}
