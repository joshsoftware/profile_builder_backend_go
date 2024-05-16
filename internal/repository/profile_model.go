package repository

type EducationDao struct {
	Degree           string `db:"degree"`
	UniversityName  string `db:"university_name"`
	Place            string `db:"place"`
	PercentageOrCgpa string `db:"percent_or_cgpa"`
	PassingYear      string `db:"passing_year"`
	CreatedAt       string `db:"created_at"`
	UpdatedAt       string `db:"updated_at"`
	CreatedById    int64  `db:"created_by_id"`
	UpdatedById    int64  `db:"updated_by_id"`
	ProfileId       int64  `db:"profile_id"`
}

type ProjectDao struct {
	Name           string `db:"name"`
	Description string `db:"description"`
	Role string `db:"role"`
	Responsibilities string `db:"responsibility"`
	Technologies string `db:"technologies"`
	TechWorkedOn string `db:"tech_worked_on"`
	WorkingStartDate string `db:"working_start_date"`
	WorkingEndDate string `db:"working_end_date"`
	Duration string `db:"duration"`
	CreatedAt       string `db:"created_at"`
	UpdatedAt       string `db:"updated_at"`
	CreatedById    int64  `db:"created_by_id"`
	UpdatedById    int64  `db:"updated_by_id"`
	ProfileId       int64  `db:"profile_id"`
}
