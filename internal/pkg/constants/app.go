package constants

import (
	"net/http"
	"time"

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
	AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodPatch},
	AllowedHeaders:   []string{"*"},
}

// CreateUserColumns defines the columns required for creating a new user profile.
var CreateUserColumns = []string{
	"name", "email", "gender", "mobile", "designation", "description", "title",
	"years_of_experience", "primary_skills", "secondary_skills", "josh_joining_date", "github_link", "linkedin_link", "career_objectives", "is_active", "is_current_employee", "created_at", "updated_at", "created_by_id", "updated_by_id",
}

// CreateEducationColumns defines the columns required for creating education details.
var CreateEducationColumns = []string{
	"degree", "university_name", "place", "percent_or_cgpa", "passing_year", "priorities", "created_at",
	"updated_at", "created_by_id", "updated_by_id", "profile_id",
}

// CreateProjectColumns defines the columns required for creating project details.
var CreateProjectColumns = []string{
	"name", "description", "role", "responsibilities", "technologies", "tech_worked_on",
	"working_start_date", "working_end_date", "duration", "priorities", "created_at", "updated_at",
	"created_by_id", "updated_by_id", "profile_id",
}

// CreateExperienceColumns defines the columns required for creating experience details.
var CreateExperienceColumns = []string{
	"designation", "company_name", "from_date", "to_date", "priorities", "created_at", "updated_at",
	"created_by_id", "updated_by_id", "profile_id",
}

// CreateCertificateColumns defines the columns required for creating certificate details.
var CreateCertificateColumns = []string{
	"name", "organization_name", "description", "issued_date", "from_date", "to_date", "priorities", "created_at", "updated_at", "created_by_id", "updated_by_id", "profile_id",
}

// CreateAchievementColumns defines the columns required for creating achievement details.
var CreateAchievementColumns = []string{
	"name", "description", "priorities", "created_at", "updated_at", "created_by_id", "updated_by_id", "profile_id",
}

// ListProfilesColumns defines the columns required for listing user profiles.
var ListProfilesColumns = []string{
	"p.id",
	"p.name",
	"p.email",
	"p.years_of_experience",
	"p.primary_skills",
	"p.is_current_employee",
	"p.is_active",
	"p.josh_joining_date",
	"p.created_at",
	"p.updated_at",
	`(SELECT 
			CASE 
				WHEN COUNT(*) = 0 THEN 0 
				WHEN COUNT(*) FILTER (WHERE is_profile_complete = 0) > 0 THEN 1 
				ELSE 0 
			END 
		FROM invitations 
		WHERE invitations.profile_id = p.id) as is_profile_complete`,
}

// ResponseProfileColumns defines the columns required for returning a specific user profile.
var ResponseProfileColumns = []string{
	"id", "name", "email", "gender", "mobile", "designation", "description", "title",
	"years_of_experience", "primary_skills", "secondary_skills", "josh_joining_date", "github_link", "linkedin_link", "career_objectives",
}

// ResponseEducationColumns defines the columns required for returning a specific user education.
var ResponseEducationColumns = []string{
	"profile_id", "id", "degree", "university_name", "place", "percent_or_cgpa", "passing_year",
}

// ResponseProjectsColumns defines the columns required for returning a specific user projects.
var ResponseProjectsColumns = []string{
	"id", "profile_id", "name", "description", "role", "responsibilities", "technologies", "tech_worked_on", "working_start_date", "working_end_date", "duration",
}

// ResponseExperiencesColumns defines the columns required for returning a specific user projects.
var ResponseExperiencesColumns = []string{
	"id", "profile_id", "designation", "company_name", "from_date", "to_date",
}

// ResponseAchievementsColumns defines the columns required for returning a specific user achievements.
var ResponseAchievementsColumns = []string{
	"id", "profile_id", "name", "description",
}

// ResponseCertificatesColumns defines the columns required for returning a specific user certificates.
var ResponseCertificatesColumns = []string{
	"id", "profile_id", "name", "organization_name", "description", "issued_date", "from_date", "to_date",
}

var RequestInvitationColumns = []string{
	"profile_id", "is_profile_complete", "created_at", "updated_at", "created_by_id", "updated_by_id",
}

// BackupTables defines the table names required for returning a backup.
var BackupTables = []string{"users", "profiles", "educations", "certificates", "projects", "experiences", "achievements"}

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

// ListQueryParams for projects
var (
	ProjectsIDsStr   = "projects_ids"
	ProjectsNamesStr = "names"
)

// ListQueryParams for educations
var (
	EducationsIDsStr   = "educations_ids"
	EducationsNamesStr = "names"
)

// ListQueryParams for experiences
var (
	ExperiencesIDsStr   = "experiences_ids"
	ExperiencesNamesStr = "names"
)

// profileID for getting query params.
var (
	ProfileID = "profile_id"
)

// ContextKey Define a custom type for context key
type ContextKey string

// Define constants for context keys
const (
	UserIDKey        ContextKey = "user_id"
	ProfileIDKey     ContextKey = "profile_id"
	AchievementIDKey ContextKey = "achievement_id"
	UserRoleKey      ContextKey = "role"
	Email            ContextKey = "email"
)

// define default values for the environment variables
var (
	// Default values for the environment variables
	DefaultMaxConnections int32         = 10
	DefaultMinConnections int32         = 0
	DefaultConnLifeTime   time.Duration = 60 * 60 // 3600 seconds
	DefaultConnIdleTime   time.Duration = 30 * 60 // 1800 seconds
	DefaultHealthCheck    time.Duration = 1 * 60  // 60 seconds
	DefaultConnectTimeout time.Duration = 5       // 5 seconds
)

// Constant Message
var (
	ResourceNotFound = "Resource not found for the given request ID"
)

// component name
var (
	Projects     = "projects"
	Achievements = "achievements"
	Educations   = "educations"
	Experiences  = "experiences"
	Certificates = "certificates"
)

// ComponentMap used to validate incoming component list
var (
	ComponentMap = map[string]bool{
		Projects:     true,
		Achievements: true,
		Educations:   true,
		Experiences:  true,
		Certificates: true,
	}
)

// DefaultMaxRetries defines the default maximum number of retries for sending an email
var (
	DefaultMaxRetries = 3
)

var (
	Admin    = "admin"
	Employee = "employee"
)

// Default profileID for the admin is 0
var (
	AdminProfileID = 0
)
