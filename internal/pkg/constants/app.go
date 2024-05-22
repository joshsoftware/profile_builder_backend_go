package constants

import (
	"net/http"

	"github.com/rs/cors"
)

const (
	EMAIL_REGEX  = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
	MOBILE_REGEX = "^([+]\\d{2})?\\d{10}$"
)

// CorsOptions defines the CORS (Cross-Origin Resource Sharing) configuration.
var CorsOptions = cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowCredentials: true,
	AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
	AllowedHeaders:   []string{"*"},
}

var CreateUserColumns = []string{"name", "email", "gender", "mobile", "designation", "description", "title", "years_of_experience", "primary_skills", "secondary_skills", "github_link", "linkedin_link", "is_active", "is_current_employee"}
