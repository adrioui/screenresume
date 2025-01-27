package models

type Departments struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DepartmentsCreate struct {
	Name string `json:"name"`
}

type DepartmentsUpdate struct {
	Name string `json:"name"`
}
