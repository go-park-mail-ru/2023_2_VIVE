package domain

type ApplicantInfo struct {
	CVid      int      `json:"cv_id"`
	FirstName string   `json:"first_name,omitempty"`
	LastName  string   `json:"last_name,omitempty"`
	Skills    []string `json:"skills,omitempty"`
}
