package domain

type DbEmployer struct {
	UserID             int    `json:"user_id"`
	OrganizationName        string `json:"organization_name"`
	OrganizationDescription string `json:"organization_description"`
}
