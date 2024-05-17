package get

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/mocks"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/stretchr/testify/mock"
)

var mockListProfile = []dto.ListProfiles{
	{
		ID:                1,
		Name:              "Abhishek Dhondalkar",
		Email:             "abhishek.dhondalkar@gmail.com",
		YearsOfExperience: 1.0,
		PrimarySkills:     []string{"Golang", "Python", "Java", "React"},
		IsCurrentEmployee: 1,
	},
}

func TestGetProfileListHandler(t *testing.T) {
	profileSvc := mocks.NewService(t)
	getProfileListHandler := GetProfileListHandler(context.Background(), profileSvc)

	tests := []struct {
		name               string
		setup              func(mock *mocks.Service)
		expectedStatusCode int
	}{
		{
			name: "Success for listing profiles",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListProfiles", mock.Anything).Return(mockListProfile, nil).Once()
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "Fail as error in ListProfiles",
			setup: func(mockSvc *mocks.Service) {
				mockSvc.On("ListProfiles", mock.Anything).Return(nil, errors.New("error")).Once()
			},
			expectedStatusCode: http.StatusBadGateway,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.setup(profileSvc)

			req, err := http.NewRequest("GET", "/list_profiles", nil)
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
