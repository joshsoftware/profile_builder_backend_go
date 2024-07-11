package test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/api/handler"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service/mocks"
	errs "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/stretchr/testify/mock"
)

func TestCreateProfileHandler(t *testing.T) {
	profileSvc := mocks.NewService(t)
	createProfileHandler := handler.CreateProfileHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		input              string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success_for_user_detail",
			input: `{ "profile" : {
                "name": "Example User",
                "email": "example.user@gmail.com",
                "gender": "Male",
                "mobile": "8888999955",
                "designation": "Employee",
                "description": "i am ml engineer",
                "title": "Software Engineer",
                "years_of_experience": 4,
				"career_objectives":"Description of the career objectives",
                "primary_skills": ["Python","SQL","Golang"],
                "secondary_skills": ["Docker", "Github"],
                "github_link": "github.com/dummy-user"
                }
            }`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("CreateProfile", mock.Anything, mock.AnythingOfType("specs.CreateProfileRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusCreated,
		},
		{
			name:               "Fail_for_incorrect_json",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			name: "Fail_for_missing_first_name_field",
			input: `
                "mobile": "9999888855",
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
			name: "Fail_for_missing_name_field",
			input: `{ "profile" : {
                "name": "",
                "email": "example.user@gmail.com",
                "gender": "Male",
                "mobile": "9999888855",
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
			name: "Fail_for_missing_email_field",
			input: `{ "profile" : {
                "name": "Example User",
                "email": "",
                "gender": "Male",
                "mobile": "9955995566",
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

var mockListProfile = []specs.ResponseListProfiles{
	{
		ID:                1,
		Name:              "Example User",
		Email:             "example@gmail.com",
		YearsOfExperience: 1.0,
		PrimarySkills:     []string{"Golang", "Python", "Java", "React"},
		IsCurrentEmployee: "YES",
	},
}

var MockSkills = []string{"GO", "RUBY", "C", "C++", "JAVA", "PYTHON", "JAVASCRIPT"}

func TestGetProfileListHandler(t *testing.T) {
	profileSvc := mocks.NewService(t)
	getProfileListHandler := handler.ProfileListHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success_for_listing_profiles",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListProfiles", mock.Anything).Return(mockListProfile, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Fail_as_error_in_listprofiles",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListProfiles", mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("GET", "/profiles", nil)
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getProfileListHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestSkillsListHandler(t *testing.T) {
	profileSvc := mocks.NewService(t)
	skillsListHandler := handler.SkillsListHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success_for_listing_skills",
			setup: func(mockSvc *mocks.Service) {
				mockListSkills := specs.ListSkills{Name: MockSkills}
				mockSvc.On("ListSkills", mock.Anything).Return(mockListSkills, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Fail_as_error_in_listskills",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListSkills", mock.Anything).Return(specs.ListSkills{}, errors.New("error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("GET", "/skills", nil)
			if err != nil {
				t.Fatal(err)
				return
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(skillsListHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestGetProfileHandler(t *testing.T) {
	profileSvc := mocks.NewService(t)
	getProfileHandler := handler.GetProfileHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		profileID          string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name:      "Success_for_getting_profile",
			profileID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetProfile", mock.Anything, 1).Return(specs.ResponseProfile{
					ProfileID:         1,
					Name:              "Example User",
					Email:             "example@example.com",
					Gender:            "Male",
					Mobile:            "1234567890",
					Designation:       "Software Engineer",
					Description:       "Experienced software engineer",
					Title:             "Full Stack Developer",
					YearsOfExperience: 5.5,
					PrimarySkills:     []string{"Go", "JavaScript"},
					SecondarySkills:   []string{"Python", "Java"},
					GithubLink:        "https://github.com/demo",
					LinkedinLink:      "https://linkedin.com/in/demo",
				}, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:      "Fail_as_error_in_getprofile",
			profileID: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("GetProfile", mock.Anything, 2).Return(specs.ResponseProfile{}, errors.New("error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("GET", "/profiles/"+test.profileID, nil)
			if err != nil {
				t.Fatal(err)
			}

			req = mux.SetURLVars(req, map[string]string{"profile_id": test.profileID})

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(getProfileHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatusCode {
				t.Errorf("Expected status code %d but got %d", test.expectedStatusCode, rr.Code)
			}
		})
	}
}

func TestUpdateProfileHandler(t *testing.T) {
	profileSvc := new(mocks.Service)
	updateProfileHandler := handler.UpdateProfileHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		url                string
		input              string
		setup              func(mockSvc *mocks.Service)
		profileID          string
		expectedStatusCode int
	}{
		{
			name:      "Success_for_updating_profile_detail",
			url:       "/profiles/1",
			profileID: "1",
			input: `{
                "profile": {
                    "id": 1,
                    "name": "Updated Name",
                    "email": "updated.email@example.com",
                    "gender": "Male",
                    "mobile": "9999999999",
                    "designation": "Senior Software Engineer",
                    "description": "Experienced software engineer with expertise in Golang",
                    "title": "Golang Developer",
                    "years_of_experience": 7,
                    "primary_skills": ["Golang", "Python"],
                    "secondary_skills": ["JavaScript", "SQL"],
                    "github_link": "https://github.com/updated",
                    "linkedin_link": "https://www.linkedin.com/in/updated"
                }
            }`,
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("UpdateProfile", context.Background(), 1, mock.AnythingOfType("specs.UpdateProfileRequest")).Return(1, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "Fail_for_rincorrect_json",
			url:                "/profiles/1",
			profileID:          "1",
			input:              "",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("PUT", test.url, bytes.NewBuffer([]byte(test.input)))
			if err != nil {
				t.Fatal(err)
				return
			}

			if test.profileID != "" {
				req = mux.SetURLVars(req, map[string]string{"profile_id": test.profileID})
			}

			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(updateProfileHandler)
			handler.ServeHTTP(rr, req)

			if rr.Result().StatusCode != test.expectedStatusCode {
				t.Errorf("Expected %d but got %d", test.expectedStatusCode, rr.Result().StatusCode)
			}
		})
	}
}

func TestDeleteProfileHandler(t *testing.T) {
	profileSvc := new(mocks.Service)

	tests := []struct {
		name               string
		profileID          string
		setup              func(mockSvc *mocks.Service)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:      "Success_for_deleting_profile",
			profileID: "1",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteProfile", mock.Anything, 1).Return(nil).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "Profile deleted successfully",
		},
		{
			name:      "No_data_found_for_deletion",
			profileID: "2",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteProfile", mock.Anything, 2).Return(errs.ErrNoData).Once()
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   "No data found for deletion",
		},
		{
			name:      "Error_while_deleting_profile",
			profileID: "3",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("DeleteProfile", mock.Anything, 3).Return(errs.ErrFailedToDelete).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "failed to delete",
		},
		{
			name:               "Error_while_getting_IDs",
			profileID:          "invalid",
			setup:              func(mockSvc *mocks.Service) {},
			expectedStatusCode: http.StatusBadGateway,
			expectedResponse:   "invalid request data",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup(profileSvc)
			reqPath := "/profiles/" + tt.profileID
			req := httptest.NewRequest(http.MethodDelete, reqPath, nil)
			req = mux.SetURLVars(req, map[string]string{"profile_id": tt.profileID})
			rr := httptest.NewRecorder()

			handler := handler.DeleteProfileHandler(context.Background(), profileSvc)
			handler(rr, req)

			res := rr.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status %d, got %d", tt.expectedStatusCode, res.StatusCode)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("could not read response: %v", err)
			}

			if !strings.Contains(string(body), tt.expectedResponse) {
				t.Errorf("expected response to contain %q, got %q", tt.expectedResponse, string(body))
			}
		})
	}
}
