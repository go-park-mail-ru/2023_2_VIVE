package domain

//easyjson:json
type ApiApplicant struct {
	CVid      int      `json:"cv_id"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Skills    []string `json:"skills,omitempty"`
}

//easyjson:json
type ApiApplicantSlice []ApiApplicant

type DbApplicant struct {
	UserID int    `json:"user_id"`
	Status string `json:"status"`
}

//easyjson:json
type ApplicantInfo struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Email     string  `json:"email"`
	AvatarURL string  `json:"avatar_url"`
	CVs       []ApiCV `json:"cvs,omitempty"`
}
