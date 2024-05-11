package request

type SignUpRequest struct {
	FistName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Gender   bool   `json:"gender"`
}
