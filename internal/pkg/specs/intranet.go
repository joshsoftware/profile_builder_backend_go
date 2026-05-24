package specs

// IntranetEmployee represents a single employee record returned by the Intranet API.
type IntranetEmployee struct {
	EmployeeID        string  `json:"employee_id"`
	Email             string  `json:"email"`
	Name              string  `json:"name"`
	MobileNumber      string  `json:"mobile_number"`
	YearsOfExperience float64 `json:"years_of_experience"`
}
