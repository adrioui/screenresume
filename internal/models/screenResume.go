package models

type CriteriaDecision struct {
	Reasoning string `json:"reasoning"`
	Decision  bool   `json:"decision"`
}

type ScreenResume struct {
	CriteriaDecisions []*CriteriaDecision `json:"criteria_decisions"`
	OverallReasoning  string              `json:"overall_reasoning"`
	OverallDecision   bool                `json:"overall_decision"`
	ResumeName        string              `json:"resume_name"`
}

type ScreenResumeCreate struct {
	ApplicationID  string   `json:"application_id"`
	JobDescription string   `json:"job_description"`
	Criteria       []string `json:"criteria"`
	File           []byte   `json:"file"`
}

type ProcessScreeningCreate struct {
	ApplicationID string   `json:"application_id"`
	Criteria      []string `json:"criteria"`
}
