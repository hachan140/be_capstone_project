package request

type AssignPeopleToManagerRequest struct {
	OrgID uint   `json:"org_id,omitempty"`
	Email string `json:"email,omitempty"`
}
