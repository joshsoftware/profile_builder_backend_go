package helpers

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
)

func SendRequest(ctx context.Context, url, access_token string) ([]byte, error) {
	serverRequest, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, errors.ErrInvalidRequestData
	}

	serverRequest.Header.Set("Authorization", "Bearer "+access_token)
	resp, err := http.DefaultClient.Do(serverRequest)
	if err != nil {
		return nil, errors.ErrHTTPRequestFailed
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.ErrReadResponseBodyFailed
	}
	return body, nil
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
