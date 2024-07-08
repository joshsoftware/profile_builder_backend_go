package helpers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
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
