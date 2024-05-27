package repository

// UserDao represents a data access object for user-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the users table.
type UserDao struct {
	ID    int64  `db:"id"`
	Email string `db:"email"`
}

// EducationDao represents a data access object for education-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the educations table.
type EducationDao struct {
	Degree           string `db:"degree"`
	UniversityName   string `db:"university_name"`
	Place            string `db:"place"`
	PercentageOrCgpa string `db:"percent_or_cgpa"`
	PassingYear      string `db:"passing_year"`
	CreatedAt        string `db:"created_at"`
	UpdatedAt        string `db:"updated_at"`
	CreatedByID      int64  `db:"created_by_id"`
	UpdatedByID      int64  `db:"updated_by_id"`
	ProfileID        int    `db:"profile_id"`
}

// ProjectDao represents a data access object for project-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the projects table.
type ProjectDao struct {
	Name             string `db:"name"`
	Description      string `db:"description"`
	Role             string `db:"role"`
	Responsibilities string `db:"responsibility"`
	Technologies     string `db:"technologies"`
	TechWorkedOn     string `db:"tech_worked_on"`
	WorkingStartDate string `db:"working_start_date"`
	WorkingEndDate   string `db:"working_end_date"`
	Duration         string `db:"duration"`
	CreatedAt        string `db:"created_at"`
	UpdatedAt        string `db:"updated_at"`
	CreatedByID      int64  `db:"created_by_id"`
	UpdatedByID      int64  `db:"updated_by_id"`
	ProfileID        int    `db:"profile_id"`
}

// ExperienceDao represents a data access object for experience-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the experiences table.
type ExperienceDao struct {
	Designation string `db:"designation"`
	CompanyName string `db:"company_name"`
	FromDate    string `db:"from_date"`
	ToDate      string `db:"to_date"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
	CreatedByID int64  `db:"created_by_id"`
	UpdatedByID int64  `db:"updated_by_id"`
	ProfileID   int    `db:"profile_id"`
}

// CertificateDao represents a data access object for certificates-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the certificates table.
type CertificateDao struct {
	Name             string `db:"name"`
	OrganizationName string `db:"organization_name"`
	Description      string `db:"description"`
	IssuedDate       string `db:"issued_date"`
	FromDate         string `db:"from_date"`
	ToDate           string `db:"to_date"`
	CreatedAt        string `db:"created_at"`
	UpdatedAt        string `db:"updated_at"`
	CreatedByID      int64  `db:"created_by_id"`
	UpdatedByID      int64  `db:"updated_by_id"`
	ProfileID        int    `db:"profile_id"`
}

// AchievementDao represents a data access object for achievements-related information.
// This struct maps to a database table, where each field corresponds to a column
// in the achievements table.
type AchievementDao struct {
	Name        string `db:"name"`
	Description string `db:"description"`
	CreatedAt   string `db:"created_at"`
	UpdatedAt   string `db:"updated_at"`
	CreatedByID int64  `db:"created_by_id"`
	UpdatedByID int64  `db:"updated_by_id"`
	ProfileID   int    `db:"profile_id"`
}
