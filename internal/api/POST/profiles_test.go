package post_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	post "github.com/joshsoftware/profile_builder_backend_go/internal/api/POST"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/mocks"
	"github.com/stretchr/testify/mock"
)

func TestCreateProfileHandler(t *testing.T) {
	profileSvc := mocks.NewService(t)
	createProfileHandler := post.CreateProfileHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for user Detail",
			input: `{ "profile" : {
                "name": "Abhishek Jain",
                "email": "abhishek.jain@gmail.com",
                "gender": "Male",
                "mobile": "9595601925",
                "designation": "Employee",
                "description": "i am ml engineer",
                "title": "Software Engineer",
                "years_of_experience": 4,
                "primary_skills": ["Python","SQL","Golang"],
                "secondary_skills": ["Docker", "Github"],
                "github_link": "github.com/dummy-user"
                }
            }`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateProfile", mock.Anything, mock.AnythingOfType("dto.CreateProfileRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Fail for incorrect json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing first_name field",
			input: `
                "mobile": "9595601925",
                "designation": "Employee",
                "description": "i am ml engineer",
                "title": "Software Engineer",
                "years_of_experience": 4,
                "primary_skills": ["Python","SQL","Golang"],
                "secondary_skills": ["Docker", "Github"],
                "github_link": "github.com/dummy-user"
                }
            }`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing name field",
			input: `{ "profile" : {
                "name": "",
                "email": "abhishek.jain@gmail.com",
                "gender": "Male",
                "mobile": "9595601925",
                "designation": "Employee",
                "description": "i am ml engineer",
                "title": "Software Engineer",
                "years_of_experience": 4,
                "primary_skills": ["Python","SQL","Golang"],
                "secondary_skills": ["Docker", "Github"],
                "github_link": "github.com/dummy-user"
                }
            }`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing email field",
			input: `{ "profile" : {
                "name": "Abhishek Jain",
                "email": "",
                "gender": "Male",
                "mobile": "9595601925",
                "designation": "Employee",
                "description": "i am ml engineer",
                "title": "Software Engineer",
                "years_of_experience": 4,
                "primary_skills": ["Python","SQL","Golang"],
                "secondary_skills": ["Docker", "Github"],
                "github_link": "github.com/dummy-user"
                }
            }`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("POST", "/profiles", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createProfileHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestCreateEducationHandler(t *testing.T) {
	userSvc := new(mocks.Service)
	createEducationHandler := post.CreateEducationHandler(context.Background(), userSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for education Detail",
			input: `{
				"profile_id": 1,
				"educations":[{
					  "degree": "BSc in Data Science",
					  "university_name": "Shivaji University",
					  "place": "Kolhapur",
					  "percent_or_cgpa": "90.50%",
					  "passing_year": "2020"
				}]
				}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateEducation", context.Background(), mock.AnythingOfType("dto.CreateEducationRequest")).Return(1, nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Fail for incorrect json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing degree field",
			input: `{
				"profile_id": 1,
				"educations":[{
					  "degree": "",
					  "university_name": "Shivaji University",
					  "place": "Kolhapur",
					  "percent_or_cgpa": "90.50%",
					  "passing_year": "2020"
				}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing profile_id field",
			input: `{
				"profile_id": 0,
				"educations":[{
					  "degree": "BSc in Data Science",
					  "university_name": "Shivaji University",
					  "place": "Kolhapur",
					  "percent_or_cgpa": "90.50%",
					  "passing_year": "2020"
				}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing passing_year field",
			input: `{
				"profile_id": 1,
				"educations":[{
					  "degree": "",
					  "university_name": "Shivaji University",
					  "place": "Kolhapur",
					  "percent_or_cgpa": "90.50%",
					  "passing_year": ""
				}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userSvc)

			req, err := http.NewRequest("POST", "/profiles/educations", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createEducationHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestCreateProjectHandler(t *testing.T) {
	userSvc := mocks.NewService(t)
	createProjectHandler := post.CreateProjectHandler(context.Background(), userSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for project Detail",
			input: `{
				"profile_id": 1,
				"projects":[{
					"name": "Least Square",
					"description": "A Webapp Which is Used to Build a Standard Profiles of an Employee for An Organization",
					"role": "Soft Developer",
					"responsibilities": "Develop a Backend",
					"technologies": "Python, Django, MongoDB, AWS",
					"tech_worked_on": "Django, AWS",
					"working_start_date": "May-2020",
					"working_end_date": "July-2020",
					"duration": "6 Months"
					}]
				}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateProject", context.Background(), mock.AnythingOfType("dto.CreateProjectRequest")).Return(1, nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Fail for incorrect json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing name field",
			input: `{
				"profile_id": 1,
				"projects":[{
					"name": "",
					"description": "A Webapp Which is Used to Build a Standard Profiles of an Employee for An Organization",
					"role": "Soft Developer",
					"responsibilities": "Develop a Backend",
					"technologies": "Python, Django, MongoDB, AWS",
					"tech_worked_on": "Django, AWS",
					"working_start_date": "May-2020",
					"working_end_date": "July-2020",
					"duration": "6 Months"
					}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing profile_id field",
			input: `{
				"profile_id": 0,
				"projects":[{
					"name": "Least Square",
					"description": "A Webapp Which is Used to Build a Standard Profiles of an Employee for An Organization",
					"role": "Soft Developer",
					"responsibilities": "Develop a Backend",
					"technologies": "Python, Django, MongoDB, AWS",
					"tech_worked_on": "Django, AWS",
					"working_start_date": "May-2020",
					"working_end_date": "July-2020",
					"duration": "6 Months"
					}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(userSvc)

			req, err := http.NewRequest("POST", "/profiles/projects", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createProjectHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestCreateExperienceHandler(t *testing.T) {
	ctx := context.Background()
	profileSvc := new(mocks.Service)
	createExperienceHandler := post.CreateExperienceHandler(ctx, profileSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for experience detail",
			input: `{
				"profile_id": 1,
				"experiences":[{
					"designation": "Associate Data Scientist",
					"company_name": "Josh Software Pvt.Ltd.",
					"from_date": "Jan-2023",
					"to_date": "July-2024"
					}]
				}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateExperience", context.Background(), mock.AnythingOfType("dto.CreateExperienceRequest")).Return(1, nil)
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Fail for incorrect JSON",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing designation field",
			input: `{
                "profile_id": 1,
                "experiences": [{
                    "designation": "",
                    "company_name": "ABC Corp",
                    "from_date": "2023-01-01",
                    "to_date": "2024-01-01"
                }]
            }`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing profile_id field",
			input: `{
                "profile_id": 0,
                "experiences": [{
                    "designation": "Software Engineer",
                    "company_name": "ABC Corp",
                    "from_date": "2023-01-01",
                    "to_date": "2024-01-01"
                }]
            }`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// mockSvc := new(mocks.Service)
			test.setup(profileSvc)

			req, err := http.NewRequest("POST", "/profiles/experiences", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createExperienceHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestCreateCertificateHandler(t *testing.T) {
	profileSvc := new(mocks.Service)
	createCertificateHandler := post.CreateCertificateHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for certificate Detail",
			input: `{
				"profile_id": 1,
				"certificates":[{
					"name": "Full Stack Data Science",
					"organization_name": "Josh Software Pvt.Ltd.",
					"description": "A Bootcamp for Mastering Data Science Concepts",
					"issued_date": "Dec-2023",
					"from_date": "June-2023",
					"to_date": "Dec-2023"
				}]
				}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateCertificate", mock.Anything, mock.AnythingOfType("dto.CreateCertificateRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Fail for incorrect json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing name field",
			input: `{
				"profile_id": 1,
				"certificates":[{
					"name": "",
					"organization_name": "Josh Software Pvt.Ltd.",
					"description": "A Bootcamp for Mastering Data Science Concepts",
					"issued_date": "Dec-2023",
					"from_date": "June-2023",
					"to_date": "Dec-2023"
				}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing profile_id field",
			input: `{
				"profile_id": 0,
				"certificates":[{
					"name": "Full Stack Data Science",
					"organization_name": "Josh Software Pvt.Ltd.",
					"description": "A Bootcamp for Mastering Data Science Concepts",
					"issued_date": "Dec-2023",
					"from_date": "June-2023",
					"to_date": "Dec-2023"
				}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing description field",
			input: `{
				"profile_id": 1,
				"certificates":[{
					"name": "Full Stack Data Science",
					"organization_name": "Josh Software Pvt.Ltd.",
					"description": "",
					"issued_date": "Dec-2023",
					"from_date": "June-2023",
					"to_date": "Dec-2023"
				}]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("POST", "/profiles/certificates", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createCertificateHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestCreateAchievementHandler(t *testing.T) {
	profileSvc := new(mocks.Service)
	createAchievementHandler := post.CreateAchievementHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for achievement detail",
			input: `{
				"profile_id": 1,
				"achievements":[{
				    "name": "Star Performer",
					"description": "Description of Award"
				    }]
				}`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateAchievement", mock.Anything, mock.AnythingOfType("dto.CreateAchievementRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Fail for incorrect json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing name field",
			input: `{
				"profile_id": 1,
				"achievements":[{
				    "name": "",
					"description": "Description of Award"
				    }]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing profile_id field",
			input: `{
				"profile_id": 0,
				"achievements":[{
				    "name": "Star Performer",
					"description": "Description of Award"
				    }]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail for missing description field",
			input: `{
				"profile_id": 1,
				"achievements":[{
				    "name": "Star Performer",
					"description": ""
				    }]
				}`,
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("POST", "/profiles/achievements", bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(createAchievementHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}
