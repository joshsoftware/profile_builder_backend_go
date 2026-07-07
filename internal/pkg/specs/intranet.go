package specs

// IntranetEmployee represents a single employee record returned by the Intranet API.
type IntranetEmployee struct {
	EmployeeID        string  `json:"employee_id"`
	Email             string  `json:"email"`
	Name              string  `json:"name"`
	MobileNumber      string  `json:"mobile_number"`
	Gender            string  `json:"gender"`
	YearsOfExperience float64 `json:"years_of_experience"`
	Designation       string  `json:"designation"`
	JoshDOJ           string  `json:"josh_doj"`
	LinkedinURL       string  `json:"linkedin_url"`
	GithubURL         string  `json:"github_url"`
	PrimarySkill      string  `json:"primary_skill"`
	SecondarySkill    string  `json:"secondary_skill"`
	Qualification     string  `json:"qualification"`
	Projects          []IntranetProject `json:"projects"`
}

// IntranetProject represents a project fetched from the Intranet API.
type IntranetProject struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"start_date"`
	EndDate     string `json:"end_date"`
}

// IntranetEmployeeResponse represents the payload returned to the frontend for form pre-fill.
type IntranetEmployeeResponse struct {
	EmployeeID        string   `json:"employeeId"`
	Email             string   `json:"email"`
	Name              string   `json:"name"`
	MobileNumber      string   `json:"mobileNumber"`
	Gender            string   `json:"gender"`
	YearsOfExperience float64  `json:"yearsOfExperience"`
	Designation       string   `json:"designation"`
	JoshJoiningDate   string   `json:"joshJoiningDate"`
	LinkedinURL       string   `json:"linkedinUrl"`
	GithubURL         string   `json:"githubUrl"`
	PrimarySkills     []string `json:"primarySkills"`
	SecondarySkills   []string `json:"secondarySkills"`
	Qualification     string   `json:"qualification"`
	Projects          []IntranetProjectResponse `json:"projects"`
}

// IntranetProjectResponse represents a project for the frontend form pre-fill.
type IntranetProjectResponse struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	StartDate   string `json:"startDate"`
	EndDate     string `json:"endDate"`
}
