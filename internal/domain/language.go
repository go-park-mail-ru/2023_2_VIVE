package domain

//easyjson:json
type Language struct {
	Name  string `json:"name"`
	Level string `json:"level,omitempty"`
}
