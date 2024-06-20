package specs

type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

type UserLoginRequest struct {
	AccessToken string `json:"access_token"`
}

type UserLoginResponse struct {
	Message    string `json:"message"`
	Token      string `json:"token"`
	StatusCode int    `json:"status_code"`
}

type UserInfo struct {
	Email string `json:"email"`
}
