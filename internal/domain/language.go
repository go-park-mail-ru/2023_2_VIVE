package domain

type Language struct {
	Name  string `json:"name"`
	Level string `json:"level,omitempty"`
}
