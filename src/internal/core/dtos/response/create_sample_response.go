package response

type CreateSampleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	StudentID   string `json:"student_id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
