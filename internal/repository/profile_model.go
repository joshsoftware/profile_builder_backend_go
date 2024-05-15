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
