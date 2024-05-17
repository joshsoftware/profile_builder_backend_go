package get

import (
	"context"
	"net/http"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/profile"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
)

func GetProfileListHandler(ctx context.Context, profileSvc profile.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		values, err := profileSvc.ListProfiles(ctx)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, dto.ListProfilesResponse{
			Profiles: values,
		})
	}
}
