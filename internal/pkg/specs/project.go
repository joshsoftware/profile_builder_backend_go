package specs

import (
	"fmt"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

// CreateProjectRequest struct represents a request to create project details.
type CreateProjectRequest struct {
	Projects []Project `json:"projects"`
}

// ListProjectsFilter used to filter projects based on ids and names
type ListProjectsFilter struct {
	ProjectsIDs []int    `json:"projects_ids"`
	Names       []string `json:"names"`
}

// Project struct represents details of a project.
type Project struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Role             string   `json:"role"`
	Responsibilities string   `json:"responsibilities"`
	Technologies     []string `json:"technologies"`
	TechWorkedOn     []string `json:"tech_worked_on"`
	WorkingStartDate string   `json:"working_start_date"`
	WorkingEndDate   string   `json:"working_end_date"`
	Duration         string   `json:"duration"`
}

// ProjectResponse struct represents details of a project for response.
type ProjectResponse struct {
	ID               int      `json:"id"`
	ProfileID        int      `json:"profile_id"`
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Role             string   `json:"role"`
	Responsibilities string   `json:"responsibilities"`
	Technologies     []string `json:"technologies"`
	TechWorkedOn     []string `json:"tech_worked_on"`
	WorkingStartDate string   `json:"working_start_date"`
	WorkingEndDate   string   `json:"working_end_date"`
	Duration         string   `json:"duration"`
}

// ResponseProject represents a project response
type ResponseProject struct {
	Projects []ProjectResponse `json:"projects"`
}

// UpdateProjectRequest represents a request to update project details.
type UpdateProjectRequest struct {
	Project Project `json:"project"`
}

// Validate func checks if the CreateProjectRequest is valid.
func (req *CreateProjectRequest) Validate() error {

	if len(req.Projects) == 0 {
		return fmt.Errorf("%s : projects ", errors.ErrEmptyPayload.Error())
	}

	for _, project := range req.Projects {
		if err := project.Validate(); err != nil {
			return err
		}
	}
	return nil
}

// Validate func checks if the Project details are valid.
func (p *Project) Validate() error {

	if p.Name == "" {
		return fmt.Errorf("%s : name ", errors.ErrParameterMissing.Error())
	}
	if p.Description == "" {
		return fmt.Errorf("%s : description ", errors.ErrParameterMissing.Error())
	}
	if p.Responsibilities == "" {
		return fmt.Errorf("%s : responsibilities ", errors.ErrParameterMissing.Error())
	}
	if len(p.TechWorkedOn) == 0 {
		return fmt.Errorf("%s : techonology worked on ", errors.ErrParameterMissing.Error())
	}
	return nil
}

// Validate func checks if the Project details are valid.
func (p *UpdateProjectRequest) Validate() error {

	fields := map[string]string{
		"name":             p.Project.Name,
		"description":      p.Project.Description,
		"responsibilities": p.Project.Responsibilities,
	}

	for fieldName, fieldValue := range fields {
		if fieldValue == "" {
			return fmt.Errorf("%s : %s", errors.ErrParameterMissing.Error(), fieldName)
		}
	}

	return nil
}

// DeleteProjectRequest represenst a Delete request of project
type DeleteProjectRequest struct {
	ProfileID int `json:"profile_id"`
	ProjectID int `json:"id"`
}
