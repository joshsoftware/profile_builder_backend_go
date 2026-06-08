package handler

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// CreateProfileHandler handles HTTP requests to create a new user profile.
func CreateProfileHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := helpers.GetUserIDFromContext(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		req, err := decodeCreateProfileRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		err = req.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		// r.Context() to send request-specific context, set by AuthMiddleware
		profileID, err := profileSvc.CreateProfile(r.Context(), req, userID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to create profile : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, specs.MessageResponseWithID{
			Message:   "Basic info added successfully",
			ProfileID: profileID,
		})
	}
}

// ProfileListHandler returns an HTTP handler that lists profiles using profileSvc.
func ProfileListHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profResponse, err := profileSvc.ListProfiles(ctx)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailedToGet)
			zap.S().Error("Unable to list profiles : ", err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.ListProfilesResponse{
			Profiles: profResponse,
		})
	}
}

// SkillsListHandler returns an HTTP handler that lists skills using profileSvc.
func SkillsListHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		skillsResp, err := profileSvc.ListSkills(ctx)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to list skills : ", err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, skillsResp)
	}
}

// GetProfileHandler returns an HTTP handler that fetches particular profile using profileSvc.
func GetProfileHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParamsByID(r, constants.ProfileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		profResp, err := profileSvc.GetProfile(ctx, profileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to get profile : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.ProfileResponse{
			Profile: profResp,
		})
	}
}

// UpdateProfileHandler returns an HTTP handler that updates profile using profileSvc.
func UpdateProfileHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParamsByID(r, constants.ProfileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		userID, err := helpers.GetUserIDFromContext(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		req, err := decodeUpdateProfileRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		err = req.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		// r.Context() to send request-specific context, set by AuthMiddleware
		updatedResp, err := profileSvc.UpdateProfile(r.Context(), profileID, userID, req)
		if err != nil {
			if err == errors.ErrAuthToken {
				middleware.ErrorResponse(w, http.StatusUnauthorized, err)
				return
			}
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailedToUpdateRecord)
			zap.S().Error("Unable to update profile : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponseWithID{
			Message:   "Basic info updated successfully",
			ProfileID: updatedResp,
		})
	}
}

// DeleteProfileHandler returns an HTTP handler that deletes profile using profileSvc.
func DeleteProfileHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParamsByID(r, constants.ProfileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("error while getting the IDs from request : ", err)
			return
		}

		err = profileSvc.DeleteProfile(ctx, profileID)
		if err != nil {
			if err == errors.ErrNoData {
				middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
					Message: constants.ResourceNotFound,
				})
				return
			}
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailedToDelete)
			zap.S().Error("error while deleting the profile: ", err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
			Message: "Profile deleted successfully",
		})
	}
}

// UpdateSequenceHandler returns an HTTP handler that updates component sequences using profileSvc.
func UpdateSequenceHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		userID, err := helpers.GetUserIDFromContext(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		req, err := decodeUpdateSequenceRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		err = req.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error(err)
			return
		}

		updatedResp, err := profileSvc.UpdateSequence(ctx, userID, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailedToUpdateRecord)
			zap.S().Error("Unable to update sequence : ", err, "for profile id : ", updatedResp)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponseWithID{
			Message:   "Component sequence updated successfully",
			ProfileID: updatedResp,
		})
	}
}

// UpdateProfileStatusHandler returns an HTTP handler that updates profile status using profileSvc.
func UpdateProfileStatusHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParamsByID(r, constants.ProfileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("error while getting the IDs from request : ", err)
			return
		}
		req, err := decodeProfileStatusRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error("error while decoding the request : ", err)
			return
		}

		err = req.ProfileStatus.Validate()
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, err)
			zap.S().Error("error while validating the request : ", err)
			return
		}

		err = profileSvc.UpdateProfileStatus(ctx, profileID, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailedToUpdateStatus)
			zap.S().Error("error while updating the profile status: ", err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
			Message: "Profile status updated successfully",
		})
	}
}

// ResolveEmployeeHandler handles request to resolve employee_id to its internal profile_id.
func ResolveEmployeeHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		employeeID, ok := vars["employee_id"]
		if !ok || employeeID == "" {
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidRequestData)
			zap.S().Error("employee_id missing from request vars")
			return
		}

		profileID, err := profileSvc.ResolveEmployeeID(ctx, employeeID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusNotFound, errors.ErrNoRecordFound)
			zap.S().Error("employee not found or resolution failed: ", err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponseWithID{
			Message:   "Employee resolved successfully",
			ProfileID: profileID,
		})
	}
}

// GetIntranetEmployeeHandler handles request to fetch employee details from Intranet for form pre-fill.
func GetIntranetEmployeeHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		employeeID, ok := vars["employee_id"]
		if !ok || employeeID == "" {
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrInvalidRequestData)
			zap.S().Error("employee_id missing from request vars")
			return
		}

		response, err := profileSvc.GetIntranetEmployee(ctx, employeeID)
		if err != nil {
			if err == errors.ErrNoRecordFound {
				middleware.ErrorResponse(w, http.StatusNotFound, errors.ErrNoRecordFound)
			} else if _, ok := err.(errors.ProfileExistsError); ok {
				middleware.ErrorResponse(w, http.StatusConflict, err)
			} else {
				middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailedToGet)
			}
			zap.S().Error("failed to get intranet employee: ", err)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, response)
	}
}
