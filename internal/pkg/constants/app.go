package constants

const (
	EmailRegex  = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
	MobileRegex = "^([+]\\d{2})?\\d{10}$"
)

var CreateUserColumns = []string{"name", "email", "gender", "mobile", "designation", "description", "title", "years_of_experience", "primary_skills", "secondary_skills", "github_link", "linkedin_link", "is_active", "is_current_employee"}

var CreateEducationColumns = []string{"degree", "university_name", "place", "percent_or_cgpa", "passing_year", "created_at", "updated_at", "created_by_id", "updated_by_id", "profile_id"}

var CreateProjectColumns = []string{"name", "description", "role", "responsibilities", "technologies", "tech_worked_on", "duration", "working_start_date", "working_end_date", "created_at", "updated_at", "created_by_id", "updated_by_id", "profile_id"}

var ListProfilesColumns = []string{"id", "name", "email", "years_of_experience", "primary_skills", "is_current_employee"}
