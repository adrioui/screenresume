package models

type JobRoles struct {
	ID           string `json:"id"`
	Title        string `json:"title"`
	DepartmentID string `json:"department_id"`
	Level        string `json:"level"`
	SalaryRange  string `json:"salary_range"`
	Location     string `json:"location"`
	IsActive     bool   `json:"is_active"`
}

type JobRolesCreate struct {
	Title        string `json:"title"`
	DepartmentID string `json:"department_id"`
	Level        string `json:"level"`
	SalaryRange  string `json:"salary_range"`
	Location     string `json:"location"`
	IsActive     bool   `json:"is_active"`
}

type JobRolesUpdate struct {
	Title        string `json:"title"`
	DepartmentID string `json:"department_id"`
	Level        string `json:"level"`
	SalaryRange  string `json:"salary_range"`
	Location     string `json:"location"`
	IsActive     bool   `json:"is_active"`
}
