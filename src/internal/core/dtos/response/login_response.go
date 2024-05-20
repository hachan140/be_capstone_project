package response

type LoginResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	Email        interface{} `json:"email,omitempty"`
}

type JWTData struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}
