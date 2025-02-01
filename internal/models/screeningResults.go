package models

import "time"

type ScreeningResults struct {
	ID            string       `json:"id"`
	ApplicationID string       `json:"application_id"`
	ModelVersion  string       `json:"model_version"`
	RawResponse   ScreenResume `json:"raw_response"`
	ProcessedAt   time.Time    `json:"processed_at"`
}

type ScreeningResultsCreate struct {
	ApplicationID string       `json:"application_id"`
	ModelVersion  string       `json:"model_version"`
	RawResponse   ScreenResume `json:"raw_response"`
}
