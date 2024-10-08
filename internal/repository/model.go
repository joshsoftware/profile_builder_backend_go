package repository

// User represents a data access object for user-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the users table.
type User struct {
	ID    int64  `db:"id"`
	Email string `db:"email"`
	Role  string `db:"role"`
	Name  string `db:"name"`
}

// ProfileRepo represents a data access object for profile-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the profile table.
type ProfileRepo struct {
	Name              string   `db:"name"`
	Email             string   `db:"email"`
	Gender            string   `db:"gender"`
	Mobile            string   `db:"mobile"`
	Designation       string   `db:"designation"`
	Description       string   `db:"description"`
	Title             string   `db:"title"`
	YearsOfExperience float64  `db:"years_of_experience"`
	PrimarySkills     []string `db:"primary_skills"`
	SecondarySkills   []string `db:"secondary_skills"`
	JoshJoiningDate   string   `db:"josh_joining_date"`
	GithubLink        string   `db:"github_link"`
	LinkedinLink      string   `db:"linkedin_link"`
	CareerObjectives  string   `db:"career_objectives"`
	CreatedAt         string   `db:"created_at"`
	UpdatedAt         string   `db:"updated_at"`
	CreatedByID       int      `db:"created_by_id"`
	UpdatedByID       int      `db:"updated_by_id"`
}

// UpdateProfileRepo represents a data access object for profile information updation.
type UpdateProfileRepo struct {
	Name              string   `db:"name"`
	Email             string   `db:"email"`
	Gender            string   `db:"gender"`
	Mobile            string   `db:"mobile"`
	Designation       string   `db:"designation"`
	Description       string   `db:"description"`
	Title             string   `db:"title"`
	YearsOfExperience float64  `db:"years_of_experience"`
	PrimarySkills     []string `db:"primary_skills"`
	SecondarySkills   []string `db:"secondary_skills"`
	JoshJoiningDate   string   `db:"josh_joining_date"`
	GithubLink        string   `db:"github_link"`
	LinkedinLink      string   `db:"linkedin_link"`
	UpdatedAt         string   `db:"updated_at"`
	UpdatedByID       int      `db:"updated_by_id"`
}

// UpdateSequenceRequest represents a data access object for component sequence updation.
type UpdateSequenceRequest struct {
	ProfileID           int         `db:"profile_id"`
	ComponentName       string      `db:"component_name"`
	ComponentPriorities map[int]int `db:"priorities"`
	UpdatedAt           string      `db:"updated_at"`
	UpdatedByID         int         `db:"updated_by_id"`
}

// EducationRepo represents a data access object for education-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the educations table.
type EducationRepo struct {
	Degree           string `db:"degree"`
	UniversityName   string `db:"university_name"`
	Place            string `db:"place"`
	PercentageOrCgpa string `db:"percent_or_cgpa"`
	PassingYear      string `db:"passing_year"`
	Priorities       int    `db:"priorities"`
	CreatedAt        string `db:"created_at"`
	UpdatedAt        string `db:"updated_at"`
	CreatedByID      int    `db:"created_by_id"`
	UpdatedByID      int    `db:"updated_by_id"`
	ProfileID        int    `db:"profile_id"`
}

// UpdateEducationRepo represents a data access object for education information updation.
type UpdateEducationRepo struct {
	Degree           string `db:"degree"`
	UniversityName   string `db:"university_name"`
	Place            string `db:"place"`
	PercentageOrCgpa string `db:"percent_or_cgpa"`
	PassingYear      string `db:"passing_year"`
	UpdatedAt        string `db:"updated_at"`
	UpdatedByID      int    `db:"updated_by_id"`
}

// ProjectRepo represents a data access object for project-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the projects table.
type ProjectRepo struct {
	Name             string   `db:"name"`
	Description      string   `db:"description"`
	Role             string   `db:"role"`
	Responsibilities string   `db:"responsibility"`
	Technologies     []string `db:"technologies"`
	TechWorkedOn     []string `db:"tech_worked_on"`
	WorkingStartDate string   `db:"working_start_date"`
	WorkingEndDate   string   `db:"working_end_date"`
	Duration         string   `db:"duration"`
	Priorities       int      `db:"priorities"`
	CreatedAt        string   `db:"created_at"`
	UpdatedAt        string   `db:"updated_at"`
	CreatedByID      int      `db:"created_by_id"`
	UpdatedByID      int      `db:"updated_by_id"`
	ProfileID        int      `db:"profile_id"`
}

// UpdateProjectRepo represents a data access object for project information updation.
type UpdateProjectRepo struct {
	Name             string   `db:"name"`
	Description      string   `db:"description"`
	Role             string   `db:"role"`
	Responsibilities string   `db:"responsibility"`
	Technologies     []string `db:"technologies"`
	TechWorkedOn     []string `db:"tech_worked_on"`
	WorkingStartDate string   `db:"working_start_date"`
	WorkingEndDate   string   `db:"working_end_date"`
	Duration         string   `db:"duration"`
	UpdatedAt        string   `db:"updated_at"`
	UpdatedByID      int      `db:"updated_by_id"`
}

// ExperienceRepo represents a data access object for experience-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the experiences table.
type ExperienceRepo struct {
	Designation string `db:"designation"`
	CompanyName string `db:"company_name"`
	FromDate    string `db:"from_date"`
	ToDate      string `db:"to_date"`
	Priorities  int    `db:"priorities"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
	CreatedByID int    `db:"created_by_id"`
	UpdatedByID int    `db:"updated_by_id"`
	ProfileID   int    `db:"profile_id"`
}

// UpdateExperienceRepo represents a data access object for experience information updation.
type UpdateExperienceRepo struct {
	Designation string `db:"designation"`
	CompanyName string `db:"company_name"`
	FromDate    string `db:"from_date"`
	ToDate      string `db:"to_date"`
	UpdatedAt   string `db:"updated_at"`
	UpdatedByID int    `db:"updated_by_id"`
}

// CertificateRepo represents a data access object for certificates-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the certificates table.
type CertificateRepo struct {
	Name             string `db:"name"`
	OrganizationName string `db:"organization_name"`
	Description      string `db:"description"`
	IssuedDate       string `db:"issued_date"`
	FromDate         string `db:"from_date"`
	ToDate           string `db:"to_date"`
	Priorities       int    `db:"priorities"`
	CreatedAt        string `db:"created_at"`
	UpdatedAt        string `db:"updated_at"`
	CreatedByID      int    `db:"created_by_id"`
	UpdatedByID      int    `db:"updated_by_id"`
	ProfileID        int    `db:"profile_id"`
}

// UpdateCertificateRepo represents a data access object for certifcate information updation.
type UpdateCertificateRepo struct {
	Name             string `db:"name"`
	OrganizationName string `db:"organization_name"`
	Description      string `db:"description"`
	IssuedDate       string `db:"issued_date"`
	FromDate         string `db:"from_date"`
	ToDate           string `db:"to_date"`
	UpdatedAt        string `db:"updated_at"`
	UpdatedByID      int    `db:"updated_by_id"`
}

// AchievementRepo represents a data access object for achievements-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the achievements table.
type AchievementRepo struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	Priorities  int    `db:"priorities"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
	CreatedByID int    `db:"created_by_id"`
	UpdatedByID int    `db:"updated_by_id"`
	ProfileID   int    `db:"profile_id"`
}

// UpdateAchievementRepo represents a data access object for achievement information updation.
type UpdateAchievementRepo struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	UpdatedAt   string `db:"updated_at"`
	UpdatedByID int    `db:"updated_by_id"`
}

// UpdateProfileStatusRepo represents a data access object for profile status information updation.
type UpdateProfileStatusRepo struct {
	IsCurrentEmployee *int   `db:"is_current_employee , omitempty"`
	IsActive          *int   `db:"is_active , omitempty"`
	UpdatedAt         string `db:"updated_at"`
}

// UpdateRequest represents a data access object for updating invitation information.
type UpdateRequest struct {
	ProfileComplete int    `db:"is_profile_complete"`
	UpdatedAt       string `db:"updated_at"`
}

// Invitations represents a data access object for email information.
type Invitations struct {
	ProfileID       int    `db:"profile_id"`
	ProfileComplete int    `db:"is_profile_complete"`
	CreatedAt       string `db:"created_at"`
	UpdatedAt       string `db:"updated_at"`
	CreatedByID     int    `db:"created_by_id"`
	UpdatedByID     int    `db:"updated_by_id"`
}

// UserInfo represents a data access object for user information.
type UserInfo struct {
	Email string `db:"email"`
	Role  string `db:"role"`
}

// GetRequest represents a data access object for getting invitation information.
type GetRequest struct {
	ProfileID         int `db:"profile_id"`
	IsProfileComplete int `db:"is_profile_complete"`
}
