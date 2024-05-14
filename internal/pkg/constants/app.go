package constants

const (
	EMAIL_REGEX  = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
	MOBILE_REGEX = "^([+]\\d{2})?\\d{10}$"
)

var CreateUserColumns = []string{"name", "email", "gender", "mobile", "designation", "description", "title", "years_of_experience", "primary_skills", "secondary_skills", "github_link", "linkedin_link", "is_active", "is_current_employee"}
