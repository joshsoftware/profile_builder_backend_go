package handler

import (
	"context"
	"net/http"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/constants"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/errors"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/specs"
	"go.uber.org/zap"
)

// CreateEducationHandler handles HTTP requests to add education details to a user profile.
func CreateEducationHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
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

		req, err := decodeCreateEducationRequest(r)
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

		profileID, err = profileSvc.CreateEducation(ctx, req, profileID, userID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to create education : ", err, " for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, specs.MessageResponseWithID{
			Message:   "Education(s) added successfully",
			ProfileID: profileID,
		})
	}
}

// ListEducationHandler returns an HTTP handler that particular education using profileSvc.
func ListEducationHandler(ctx context.Context, eduSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParamsByID(r, constants.ProfileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}
		filter, err := helpers.DecodeEducationsRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrDecodeRequest)
			zap.S().Error(err)
			return
		}

		eduResp, err := eduSvc.ListEducations(ctx, profileID, filter)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailespecsFetch)
			zap.S().Error("Unable to fetch educations : ", err, "for profile id : ", profileID)
			return
		}

		if len(eduResp) == 0 {
			eduResp = []specs.EducationResponse{}
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.ResponseEducation{
			Educations: eduResp,
		})
	}
}

// UpdateEducationHandler returns an HTTP handler that updates education using profileSvc.
func UpdateEducationHandler(ctx context.Context, eduSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, eduID, err := helpers.GetMultipleParams(r)
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

		req, err := decodeUpdateEducationRequest(r)
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

		updatedResp, err := eduSvc.UpdateEducation(ctx, profileID, eduID, userID, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to update education : ", err, "for profile id : ", profileID, "education id : ", eduID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponseWithID{
			Message:   "Education updated successfully",
			ProfileID: updatedResp,
		})
	}
}
