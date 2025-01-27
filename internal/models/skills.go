package models

type Skills struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Category string `json:"category"`
}

type SkillsCreate struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

type SkillsUpdate struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}
