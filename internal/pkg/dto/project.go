package dto

import (
	"errors"
)

type CreateProjectRequest struct {
	ProfileId int64       `json:"profile_id"`
	Projects []Project `json:"projects"`
}

type Project struct {
	Name           string `json:"name"`
	Description string `json:"description"`
	Role string `json:"role"`
	Responsibilities string `json:"responsibilities"`
	Technologies string `json:"technologies"`
	TechWorkedOn string `json:"tech_worked_on"`
	WorkingStartDate string `json:"working_start_date"`
	WorkingEndDate string `json:"working_end_date"`
	Duration string `json:"duration"`
}

func (req *CreateProjectRequest) Validate() error {
	
	if req.ProfileId <= 0 {
		return errors.New("profile_id must be a positive integer")
	}
	
	for _, project := range req.Projects {
		if err := project.Validate(); err != nil {
			return err
		}
	}
	return nil
}

func (p *Project) Validate() error {

	if p.Name == "" {
		return errors.New("name is required")
	}
	if p.Description == "" {
		return errors.New("description is required")
	}
	if p.Role == "" {
		return errors.New("role is required")
	}
	if p.Responsibilities == "" {
		return errors.New("responsibility is required")
	}
	if p.Technologies == "" {
		return errors.New("technologies is required")
	}
	if p.TechWorkedOn == "" {
		return errors.New("tech_worked_on is required")
	}
	if p.WorkingStartDate == "" {
		return errors.New("working_start_date is required")
	}
	if p.WorkingEndDate == "" {
		return errors.New("working_end_date is required")
	}
	if p.Duration == "" {
		return errors.New("duration is required")
	}

	return nil
}
