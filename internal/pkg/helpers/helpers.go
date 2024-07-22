package helpers

import (
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// GetQueryIntIds used to get query params as in int format
func GetQueryIntIds(r *http.Request, key string) ([]int, error) {
	var IDInts []int
	if key != "" {
		IDStrs := strings.Split(key, ",")
		for _, idStr := range IDStrs {
			idInt, err := strconv.Atoi(idStr)
			if err != nil {
				return IDInts, errors.ErrInvalidRequestData
			}
			IDInts = append(IDInts, idInt)
		}
	}
	return IDInts, nil
}

// GetQueryStrings used to get query params as in strings format
func GetQueryStrings(r *http.Request, key string) []string {
	var NameStrs []string
	if key != "" {
		NameStrs = strings.Split(key, ",")
	}
	return NameStrs
}

// ConvertStringToIntWithDefault returns the integer value of given string
func ConvertStringToIntWithDefault(envVars string, defaultValue int32) int32 {
	valueStr := os.Getenv(envVars)
	if valueStr == "" {
		zap.S().Info("Environment variable is not provided ,setting default value ", defaultValue, " for key : ", envVars)
		return defaultValue
	}

	value, err := strconv.ParseInt(valueStr, 10, 32)
	if err != nil {
		zap.S().Warn("convestion failes for env variable ", envVars, "expect string but got : ", reflect.TypeOf(valueStr), " : ", err)
		return defaultValue
	}
	return int32(value)
}

// ConvertStringToTimeDuration converts a String to time duration
func ConvertStringToTimeDuration(envVars string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(envVars)
	if valueStr == "" {
		zap.S().Info("Environment variable is not provided ,setting default value ", defaultValue, " for key : ", envVars)
		return defaultValue
	}

	value, err := time.ParseDuration(valueStr)
	if err != nil {
		zap.S().Warn("convestion failes for env variable ", envVars, "expect string but got : ", reflect.TypeOf(valueStr), " : ", err)
		return defaultValue
	}
	return value
}

// constructEmailMessage constructs the email message
func ConstructEmailMessage(from string, request specs.UserEmailRequest) []byte {
	// Construct the link using the profile ID
	link := fmt.Sprintf("http://localhost:3000/profile-builder/%d", request.ProfileID)

	messageBody := fmt.Sprintf(`
        <html>
        <head></head>
        <body>
            <p>Your profile has been created successfully. Please click on the link below to see the details:</p>
            <p><a href="%s">Click here</a></p>
        </body>
        </html>`, link)

	return []byte(strings.Join([]string{
		"From: " + from,
		"To: " + request.Email,
		"Subject: Update your profile",
		"MIME-version: 1.0;",
		"Content-Type: text/html; charset=\"UTF-8\";",
		"",
		messageBody,
	}, "\r\n"))
}
