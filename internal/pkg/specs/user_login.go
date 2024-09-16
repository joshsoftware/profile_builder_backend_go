package specs

// UserInfoFilter struct to store user details
type UserInfoFilter struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

// UserLoginRequest to get token detail
type UserLoginRequest struct {
	AccessToken string `json:"access_token"`
}

// UserLoginResponse to respond with login
type UserLoginResponse struct {
	Message    string `json:"message"`
	ProfileID  int    `json:"profile_id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	Token      string `json:"token"`
	StatusCode int    `json:"status_code"`
}

// LoginResponse to respond with login
type LoginResponse struct {
	ProfileID int    `json:"profile_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Token     string `json:"token"`
}
