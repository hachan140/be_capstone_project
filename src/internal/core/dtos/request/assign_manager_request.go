package request

type AssignPeopleToManagerRequest struct {
	Email string `json:"email,omitempty"`
}

type RecallPeopleManagerRequest struct {
	Email string `json:"email"`
}
