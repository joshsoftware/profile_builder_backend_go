package constants

import (
	"net/http"

	"github.com/rs/cors"
)

// Email and Mobile Regex defines a regular expression pattern for validating email addresses
// and mobile
const (
	EmailRegex  = "^[\\w-\\.]+@([\\w-]+\\.)+[\\w-]{2,4}$"
	MobileRegex = "^([+]\\d{2})?\\d{10}$"
)

// CorsOptions defines the CORS (Cross-Origin Resource Sharing) configuration.
var CorsOptions = cors.Options{
	AllowedOrigins:   []string{"*"},
	AllowCredentials: true,
	AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
	AllowedHeaders:   []string{"*"},
}

// CreateUserColumns defines the columns required for creating a new user profile.
var CreateUserColumns = []string{
	"name", "email", "gender", "mobile", "designation", "description", "title",
	"years_of_experience", "primary_skills", "secondary_skills", "github_link", "linkedin_link",
	"is_active", "is_current_employee", "created_at", "updated_at", "created_by_id", "updated_by_id",
}

// CreateEducationColumns defines the columns required for creating education details.
var CreateEducationColumns = []string{
	"degree", "university_name", "place", "percent_or_cgpa", "passing_year", "created_at",
	"updated_at", "created_by_id", "updated_by_id", "profile_id",
}

// CreateProjectColumns defines the columns required for creating project details.
var CreateProjectColumns = []string{
	"name", "description", "role", "responsibilities", "technologies", "tech_worked_on",
	"duration", "working_start_date", "working_end_date", "created_at", "updated_at",
	"created_by_id", "updated_by_id", "profile_id",
}

// CreateExperienceColumns defines the columns required for creating experience details.
var CreateExperienceColumns = []string{
	"designation", "company_name", "from_date", "to_date", "created_at", "updated_at",
	"created_by_id", "updated_by_id", "profile_id",
}

// CreateCertificateColumns defines the columns required for creating certificate details.
var CreateCertificateColumns = []string{
	"name", "organization_name", "description", "issued_date", "from_date", "to_date", "created_at", "updated_at", "created_by_id", "updated_by_id", "profile_id",
}

// CreateAchievementColumns defines the columns required for creating achievement details.
var CreateAchievementColumns = []string{
	"name", "description", "created_at", "updated_at", "created_by_id", "updated_by_id", "profile_id",
}

// ListProfilesColumns defines the columns required for listing user profiles.
var ListProfilesColumns = []string{
	"id", "name", "email", "years_of_experience", "primary_skills", "is_current_employee",
}

// ResponseProfileColumns defines the columns required for returning a specific user profile.
var ResponseProfileColumns = []string{
	"id", "name", "email", "gender", "mobile", "designation", "description", "title",
	"years_of_experience", "primary_skills", "secondary_skills", "github_link", "linkedin_link",
}

// ResponseEducationColumns defines the columns required for returning a specific user education.
var ResponseEducationColumns = []string{
	"profile_id", "degree", "university_name", "place", "percent_or_cgpa", "passing_year",
}

// ResponseProjectsColumns defines the columns required for returning a specific user projects.
var ResponseProjectsColumns = []string{
	"profile_id", "name", "description", "role", "responsibilities", "technologies", "tech_worked_on",
	"duration", "working_start_date", "working_end_date",
}

// ResponseExperiencesColumns defines the columns required for returning a specific user projects.
var ResponseExperiencesColumns = []string{
	"profile_id", "designation", "company_name", "from_date", "to_date",
}

// ResponseAchievementsColumns defines the columns required for returning a specific user achievements.
var ResponseAchievementsColumns = []string{
	"id", "profile_id", "name", "description",
}

// ResponseCertificatesColumns defines the columns required for returning a specific user certificates.
var ResponseCertificatesColumns = []string{
	"id", "profile_id", "name", "organization_name", "description", "issued_date", "from_date", "to_date",
}

// ListQueryParams for acheivements
var (
	AchievementIDsStr   = "achievement_ids"
	AchievementNamesStr = "names"
)

// ListQueryParams for certificates
var (
	CertificateIDsStr   = "certificate_ids"
	CertificateNamesStr = "names"
)
