package domain

type ApiApplicant struct {
	CVid      int      `json:"cv_id"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Skills    []string `json:"skills,omitempty"`
}

type DbApplicant struct {
	UserID int    `json:"user_id"`
	Status string `json:"status"`
}

type ApplicantInfo struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	CVs       []ApiCV `json:"cvs,omitempty"`
}
