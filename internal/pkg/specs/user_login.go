package specs

// User to store ID and Email from user input
type User struct {
	ID    int64  `json:"id"`
	Email string `json:"email"`
}

// UserLoginRequest to get token detail
type UserLoginRequest struct {
	AccessToken string `json:"access_token"`
}

// UserLoginResponse to respond with login
type UserLoginResponse struct {
	Message    string `json:"message"`
	Token      string `json:"token"`
	StatusCode int    `json:"status_code"`
}

// UserInfo for getting mail of specified login
type UserInfo struct {
	Email string `json:"email"`
}
