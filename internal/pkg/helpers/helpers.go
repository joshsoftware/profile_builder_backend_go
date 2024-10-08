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

// JoinValues used to join multiple values while data backing up
func JoinValues(values interface{}, sep string) string {
	switch v := values.(type) {
	case []string:
		return strings.Join(v, sep)
	case []interface{}:
		strValues := make([]string, len(v))
		for i, val := range v {
			strValues[i] = fmt.Sprintf("%v", val)
		}
		return strings.Join(strValues, sep)
	default:
		return ""
	}
}

// ConstructUserMessage constructs the email message for a profile invitation
func ConstructUserMessage(email, name string, profileID int) string {
	link := fmt.Sprintf("%s/profile-builder/%d", os.Getenv("HOST_URL"), profileID)
	content := fmt.Sprintf(`
		<html>
		<body>
			<div class="email-content">
				<p>Hello %s,</p>
				<p>We are pleased to inform you that your Josh profile has been successfully created in Profile Builder.</p>
				<p>Please <a href="%s">click here</a> to review your profile and update the remaining details as soon as possible.</p>
				<p>Once all the required information has been provided, kindly submit your profile for final approval.</p>
				<p>Feel free to reach out to Ruchira Kulkarni if you have any questions or need assistance.</p>
				<p>Best Regards,</p>
				<p>Profile Builder Team</p>
			</div>
		</body>
		</html>
	`, name, link)
	return content
}

// ConstructAdminEmailMessage constructs the email message for an admin invitation
func ConstructAdminEmailMessage(email, name string, profileID int) string {
	link := fmt.Sprintf("%s/profile-builder/%d", os.Getenv("HOST_URL"), profileID)
	content := fmt.Sprintf(`
		<html>
		<body>
			<div class="email-content">
				<p>Hello %s,</p>
				<p>The candidate has completed their profile. Please <a href="%s">click here</a> to review and download the profile.</p>
				<p>Best Regards,</p>
				<p>Profile Builder Team</p>
			</div>
		</body>
		</html>
	`, name, link)
	return content
}

// GetCurrentISTTime returns the current time in the Asia/Kolkata time zone formatted as RFC3339
func GetCurrentISTTime() string {
	loc, err := time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return ""
	}
	return time.Now().In(loc).Format(time.RFC3339)
}

// CheckBoolStatus returns YES if value is 1 else NO
func CheckBoolStatus(value int) string {
	if value == 1 {
		return "YES"
	}
	return "NO"
}

// ConvertToLowerCase converts the email to lowercase
func ConvertToLowerCase(email string) string {
	return strings.ToLower(email)
}
