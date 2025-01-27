package models

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
