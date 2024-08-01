package helpers

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.uber.org/zap"
)

// SendRequest used to check valid login request
func SendRequest(ctx context.Context, methodType, url, accessToken string, body io.Reader, headers map[string]string) ([]byte, error) {
	serverRequest, err := http.NewRequestWithContext(ctx, methodType, url, body)
	if err != nil {
		log.Fatalf("Error in creating request: %v", err)
		return nil, err
	}

	if accessToken != "" {
		serverRequest.Header.Set("Authorization", "Bearer "+accessToken)
	}

	for key, value := range headers {
		serverRequest.Header.Set(key, value)
	}

	resp, err := http.DefaultClient.Do(serverRequest)
	if err != nil {
		return nil, errors.ErrHTTPRequestFailed
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.ErrReadResponseBodyFailed
	}

	return respBody, nil
}

// IsDuplicateKeyError returns true if the given key is duplicate and false otherwise
func IsDuplicateKeyError(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return strings.Contains(pgErr.Message, "duplicate key value violates unique constraint")
	}
	return false
}

// IsInvalidProfileError returns true if the given profile_id is wrong and false otherwise
func IsInvalidProfileError(err error) bool {
	if pgErr, ok := err.(*pgconn.PgError); ok {
		return strings.Contains(pgErr.Message, "violates foreign key constraint")
	}
	return false
}

// GetTodaysDate returns the current date in string format
func GetTodaysDate() string {
	now := time.Now()
	today := strings.Split(now.String(), " ")
	return today[0]
}

// ConvertStringToInt returns the integer value of given string
func ConvertStringToInt(value string) (int, error) {
	id, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.ErrInvalidRequestData
	}
	return id, nil
}

// GetParamsByID returns the url params which is coming from the query parameters
func GetParamsByID(r *http.Request, id string) (ID int, err error) {
	vars := mux.Vars(r)
	id, ok := vars[id]
	if !ok {
		return 0, errors.ErrInvalidRequestData
	}

	ID, err = strconv.Atoi(id)
	if err != nil {
		return 0, errors.ErrInvalidRequestData
	}
	return ID, nil
}

// GetUserIDFromContext returns the user_id which is coming from the context
func GetUserIDFromContext(r *http.Request) (ID int, err error) {
	userID, ok := r.Context().Value(constants.UserIDKey).(float64)
	if !ok {
		return 0, errors.ErrInvalidUserID
	}
	return int(userID), nil
}

func GetContextValue(r *http.Request) (specs.UserContext, error) {
	email, ok := r.Context().Value(constants.Email).(string)
	if !ok {
		return specs.UserContext{}, errors.ErrInvalidEmail
	}

	role, ok := r.Context().Value(constants.UserRoleKey).(string)
	if !ok {
		return specs.UserContext{}, errors.ErrUserRole
	}

	userContext := specs.UserContext{
		Role:  role,
		Email: email,
	}
	return userContext, nil
}

// GetMultipleParams returns the multiple IDs which is coming from the query parameters
func GetMultipleParams(r *http.Request) (int, int, error) {
	vars := mux.Vars(r)

	profileIDStr, profileIDOk := vars["profile_id"]
	idStr, idOk := vars["id"]

	if !profileIDOk || !idOk {
		return 0, 0, errors.ErrInvalidRequestData
	}

	profileID, err := strconv.Atoi(profileIDStr)
	if err != nil {
		return 0, 0, errors.ErrInvalidRequestData
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, 0, errors.ErrInvalidRequestData
	}

	return profileID, id, nil
}

// DecodeAchievementRequest Decode the ListFilter
func DecodeAchievementRequest(r *http.Request) (specs.ListAchievementFilter, error) {
	achievementIDs := r.URL.Query().Get(constants.AchievementIDsStr)
	achievementNames := r.URL.Query().Get(constants.AchievementNamesStr)

	if achievementIDs == "" && achievementNames == "" {
		return specs.ListAchievementFilter{}, nil
	}

	achievementIntIDs, err := GetQueryIntIds(r, achievementIDs)
	if err != nil {
		return specs.ListAchievementFilter{}, err
	}

	achievementNamesStr := GetQueryStrings(r, achievementNames)
	filter := specs.ListAchievementFilter{
		AchievementIDs: achievementIntIDs,
		Names:          achievementNamesStr,
	}
	return filter, nil
}

// DecodeCertificateRequest decode Certificate request and returns a filter
func DecodeCertificateRequest(r *http.Request) (specs.ListCertificateFilter, error) {
	certificateIDs := r.URL.Query().Get(constants.CertificateIDsStr)
	certificateNames := r.URL.Query().Get(constants.CertificateNamesStr)

	if certificateIDs == "" && certificateNames == "" {
		return specs.ListCertificateFilter{}, nil
	}

	certificateIDsInts, err := GetQueryIntIds(r, certificateIDs)
	if err != nil {
		return specs.ListCertificateFilter{}, err
	}

	certificateNameStrs := GetQueryStrings(r, certificateNames)

	filter := specs.ListCertificateFilter{
		CertificateIDs: certificateIDsInts,
		Names:          certificateNameStrs,
	}
	return filter, nil
}

// DecodeProjectsRequest decode Projects request and returns a filter
func DecodeProjectsRequest(r *http.Request) (specs.ListProjectsFilter, error) {
	projectsIDs := r.URL.Query().Get(constants.ProjectsIDsStr)
	projectsNames := r.URL.Query().Get(constants.ProjectsNamesStr)

	if projectsIDs == "" && projectsNames == "" {
		return specs.ListProjectsFilter{}, nil
	}

	projectsIDsInts, err := GetQueryIntIds(r, projectsIDs)
	if err != nil {
		return specs.ListProjectsFilter{}, err
	}

	projectsNameStrs := GetQueryStrings(r, projectsNames)

	filter := specs.ListProjectsFilter{
		ProjectsIDs: projectsIDsInts,
		Names:       projectsNameStrs,
	}
	return filter, nil
}

// DecodeEducationsRequest decode educations request and returns a filter
func DecodeEducationsRequest(r *http.Request) (specs.ListEducationsFilter, error) {
	educationsIDs := r.URL.Query().Get(constants.EducationsIDsStr)
	educationsNames := r.URL.Query().Get(constants.EducationsNamesStr)

	if educationsIDs == "" && educationsNames == "" {
		return specs.ListEducationsFilter{}, nil
	}

	projectsIDsInts, err := GetQueryIntIds(r, educationsIDs)
	if err != nil {
		return specs.ListEducationsFilter{}, err
	}

	projectsNameStrs := GetQueryStrings(r, educationsNames)

	filter := specs.ListEducationsFilter{
		EduationsIDs: projectsIDsInts,
		Names:        projectsNameStrs,
	}
	return filter, nil
}

// DecodeExperiencesRequest decode experiences request and returns a filter
func DecodeExperiencesRequest(r *http.Request) (specs.ListExperiencesFilter, error) {
	experiencesIDs := r.URL.Query().Get(constants.ExperiencesIDsStr)
	experiencesNames := r.URL.Query().Get(constants.ExperiencesNamesStr)

	if experiencesIDs == "" && experiencesNames == "" {
		return specs.ListExperiencesFilter{}, nil
	}

	experiencesIDsInts, err := GetQueryIntIds(r, experiencesIDs)
	if err != nil {
		return specs.ListExperiencesFilter{}, err
	}

	experiencesNameStrs := GetQueryStrings(r, experiencesNames)

	filter := specs.ListExperiencesFilter{
		ExperiencesIDs: experiencesIDsInts,
		Names:          experiencesNameStrs,
	}
	return filter, nil
}

// getEmailConfig returns the email configuration from the environment variables
func getEmailConfig() (from, api_key string) {
	from = os.Getenv("FROM_EMAIL")
	api_key = os.Getenv("SENDGRID_API_KEY")
	return
}

// SendAdminInvitation sends an admin invitation email
func SendAdminInvitation(email string, profileID int) error {
	message := ConstructAdminEmailMessage(email, profileID)
	return SendInvitation(email, "Admin Invitation", message)
}

// SendUserInvitation sends a user invitation email
func SendUserInvitation(email string, profileID int) error {
	message := ConstructUserMessage(email, profileID)
	return SendInvitation(email, "Profile Invitation", message)
}

// SendInvitation sends an invitation email
func SendInvitation(email string, subject string, message string) error {
	from, sendgridAPIKey := getEmailConfig()
	client := sendgrid.NewSendClient(sendgridAPIKey)

	fromEmail := mail.NewEmail("Saurabh Puri", from)
	toEmail := mail.NewEmail("Saurabh Puri", email)
	plainTextContent := message
	htmlContent := message
	msg := mail.NewSingleEmail(fromEmail, subject, toEmail, plainTextContent, htmlContent)

	response, err := client.Send(msg)
	if err != nil {
		zap.S().Error("failed to send email : ", err)
		return err
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		zap.S().Info("Email sent successfully", zap.String("email", email))
	} else {
		zap.S().Error("Failed to send email", zap.String("email", email), zap.String("response", response.Body))
	}
	return nil
}

func GetProfileId(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	profileIDStr, ok := vars["profile_id"]
	if !ok {
		return 0, errors.ErrInvalidRequestData
	}

	profileID, err := strconv.Atoi(profileIDStr)
	if err != nil {
		return 0, errors.ErrInvalidRequestData
	}
	return profileID, nil
}
