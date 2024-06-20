package helpers

import (
	"context"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

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

// MultipleConvertStringToInt returns the integer value of given string
func MultipleConvertStringToInt(profileID string, id string) (int, int, error) {
	profileIDInt, err := strconv.Atoi(profileID)
	if err != nil {
		return 0, 0, errors.ErrInvalidRequestData
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return 0, 0, errors.ErrInvalidRequestData
	}

	return profileIDInt, idInt, nil
}

// GetParams returns the profileID which is coming from the query parameters
func GetParams(r *http.Request) (ID string, err error) {
	vars := mux.Vars(r)
	profileID, ok := vars["profile_id"]
	if !ok {
		return "", errors.ErrInvalidRequestData
	}
	return profileID, nil
}

// GetMultipleParams returns the multiple IDs which is coming from the query parameters
func GetMultipleParams(r *http.Request) (string, string, error) {
	vars := mux.Vars(r)

	profileID, profileIDOk := vars["profile_id"]
	id, idOk := vars["id"]

	if !profileIDOk || !idOk {
		return "", "", errors.ErrInvalidRequestData
	}

	return profileID, id, nil
}

// Decode the ListFilter
func DecodeAchievementRequest(r *http.Request) (dto.ListAchievementFilter, error) {
	achievementIDs := r.URL.Query().Get(constants.AchievementIDsStr)
	achievementNames := r.URL.Query().Get(constants.AchievementNamesStr)

	if achievementIDs == "" && achievementNames == "" {
		return dto.ListAchievementFilter{}, nil
	}

	achievementIntIDs, err := GetQueryIntIds(r, achievementIDs)
	if err != nil {
		return dto.ListAchievementFilter{}, err
	}

	achievementNamesStr := GetQueryStrings(r, achievementNames)
	filter := dto.ListAchievementFilter{
		AchievementIDs: achievementIntIDs,
		Names:          achievementNamesStr,
	}
	return filter, nil
}

func DecodeCertificateRequest(r *http.Request) (dto.ListCertificateFilter, error) {
	certificateIDs := r.URL.Query().Get(constants.CertificateIDsStr)
	certificateNames := r.URL.Query().Get(constants.CertificateNamesStr)

	if certificateIDs == "" && certificateNames == "" {
		return dto.ListCertificateFilter{}, nil
	}

	certificateIDsInts, err := GetQueryIntIds(r, certificateIDs)
	if err != nil {
		return dto.ListCertificateFilter{}, err
	}

	certificateNameStrs := GetQueryStrings(r, certificateNames)

	filter := dto.ListCertificateFilter{
		CertificateIDs: certificateIDsInts,
		Names:          certificateNameStrs,
	}
	return filter, nil
}
