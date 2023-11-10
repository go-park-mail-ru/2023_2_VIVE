package domain

type Organization struct {
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Location    *string `json:"location,omitempty"`
}
