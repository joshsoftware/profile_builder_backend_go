package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	clientmocks "github.com/joshsoftware/profile_builder_backend_go/internal/client/intranet/mocks"
	pkgerrors "github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	repomocks "github.com/joshsoftware/profile_builder_backend_go/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSyncEmployees(t *testing.T) {
	tests := []struct {
		name            string
		setup           func(profileMock *repomocks.ProfileStorer, intranetMock *clientmocks.IntranetClient)
		wantUpdated     int
		wantSkipped     int
		isErrorExpected bool
	}{
		{
			name: "Success_all_matched",
			setup: func(profileMock *repomocks.ProfileStorer, intranetMock *clientmocks.IntranetClient) {
				employees := []specs.IntranetEmployee{
					{EmployeeID: "EMP001", Email: "alice@example.com"},
					{EmployeeID: "EMP002", Email: "bob@example.com"},
				}
				intranetMock.On("GetEmployees", mock.Anything).Return(employees, nil).Once()

				profileMock.On("UpdateEmployeeIDByEmail", mock.Anything, "alice@example.com", "EMP001").Return(nil).Once()
				profileMock.On("UpdateEmployeeIDByEmail", mock.Anything, "bob@example.com", "EMP002").Return(nil).Once()
			},
			wantUpdated:     2,
			wantSkipped:     0,
			isErrorExpected: false,
		},
		{
			name: "Partial_match",
			setup: func(profileMock *repomocks.ProfileStorer, intranetMock *clientmocks.IntranetClient) {
				employees := []specs.IntranetEmployee{
					{EmployeeID: "EMP001", Email: "alice@example.com"},
					{EmployeeID: "EMP002", Email: "bob@example.com"},
					{EmployeeID: "EMP003", Email: "unknown@example.com"},
				}
				intranetMock.On("GetEmployees", mock.Anything).Return(employees, nil).Once()

				profileMock.On("UpdateEmployeeIDByEmail", mock.Anything, "alice@example.com", "EMP001").Return(nil).Once()
				profileMock.On("UpdateEmployeeIDByEmail", mock.Anything, "bob@example.com", "EMP002").Return(nil).Once()
				profileMock.On("UpdateEmployeeIDByEmail", mock.Anything, "unknown@example.com", "EMP003").Return(pkgerrors.ErrNoRecordFound).Once()
			},
			wantUpdated:     2,
			wantSkipped:     1,
			isErrorExpected: false,
		},
		{
			name: "No_match",
			setup: func(profileMock *repomocks.ProfileStorer, intranetMock *clientmocks.IntranetClient) {
				employees := []specs.IntranetEmployee{
					{EmployeeID: "EMP099", Email: "ghost1@example.com"},
					{EmployeeID: "EMP100", Email: "ghost2@example.com"},
				}
				intranetMock.On("GetEmployees", mock.Anything).Return(employees, nil).Once()

				profileMock.On("UpdateEmployeeIDByEmail", mock.Anything, "ghost1@example.com", "EMP099").Return(pkgerrors.ErrNoRecordFound).Once()
				profileMock.On("UpdateEmployeeIDByEmail", mock.Anything, "ghost2@example.com", "EMP100").Return(pkgerrors.ErrNoRecordFound).Once()
			},
			wantUpdated:     0,
			wantSkipped:     2,
			isErrorExpected: false,
		},
		{
			name: "IntranetClient_error",
			setup: func(profileMock *repomocks.ProfileStorer, intranetMock *clientmocks.IntranetClient) {
				intranetMock.On("GetEmployees", mock.Anything).Return(nil, errors.New("API unavailable")).Once()
			},
			wantUpdated:     0,
			wantSkipped:     0,
			isErrorExpected: true,
		},
		{
			name: "UpdateEmployeeIDByEmail_error_one_employee_continues",
			setup: func(profileMock *repomocks.ProfileStorer, intranetMock *clientmocks.IntranetClient) {
				employees := []specs.IntranetEmployee{
					{EmployeeID: "EMP001", Email: "alice@example.com"},
					{EmployeeID: "EMP002", Email: "bob@example.com"},
				}
				intranetMock.On("GetEmployees", mock.Anything).Return(employees, nil).Once()

				// alice — DB error (non-fatal: log and continue)
				profileMock.On("UpdateEmployeeIDByEmail", mock.Anything, "alice@example.com", "EMP001").Return(errors.New("db error")).Once()
				// bob — succeeds
				profileMock.On("UpdateEmployeeIDByEmail", mock.Anything, "bob@example.com", "EMP002").Return(nil).Once()
			},
			wantUpdated:     1,
			wantSkipped:     1,
			isErrorExpected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProfileRepo := new(repomocks.ProfileStorer)
			mockIntranetClient := new(clientmocks.IntranetClient)

			repoDeps := service.RepoDeps{
				ProfileDeps:    mockProfileRepo,
				IntranetClient: mockIntranetClient,
			}
			svc := service.NewServices(repoDeps)

			tt.setup(mockProfileRepo, mockIntranetClient)

			updated, skipped, err := svc.SyncEmployees(context.Background())

			assert.Equal(t, tt.wantUpdated, updated)
			assert.Equal(t, tt.wantSkipped, skipped)
			if tt.isErrorExpected {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockProfileRepo.AssertExpectations(t)
			mockIntranetClient.AssertExpectations(t)
		})
	}
}
