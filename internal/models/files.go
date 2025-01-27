package models

type Files struct {
	ID       string `json:"id"`
	Path     string `json:"path"`
	FileType string `json:"file_type"`
	Checksum string `json:"checksum"`
}

type FilesCreate struct {
	Path     string `json:"path"`
	FileType string `json:"file_type"`
	Checksum string `json:"checksum"`
}

type FilesUpdate struct {
	Path     string `json:"path"`
	FileType string `json:"file_type"`
	Checksum string `json:"checksum"`
}
