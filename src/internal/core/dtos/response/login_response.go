package response

type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	Email       interface{} `json:"email"`
}

type JWTData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
