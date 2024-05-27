package handler

import (
	"context"
	"net/http"

	"github.com/joshsoftware/profile_builder_backend_go/internal/app/service"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/dto"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/helpers"
	"github.com/joshsoftware/profile_builder_backend_go/internal/pkg/middleware"
	"go.uber.org/zap"
)

// CreateEducationHandler handles HTTP requests to add education details to a user profile.
func CreateEducationHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		profileID, err := profileSvc.CreateEducation(ctx, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to create education : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, dto.MessageResponseWithID{
			Message:   "Education(s) added successfully",
			ProfileID: profileID,
		})
	}
}

// GetEducationHandler returns an HTTP handler that particular education using profileSvc.
func GetEducationHandler(ctx context.Context, eduSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParams(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		values, err := eduSvc.GetEducation(ctx, profileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to fetch education : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, dto.ResponseEducation{
			Educations: values,
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

		value, err := eduSvc.UpdateEducation(ctx, profileID, eduID, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to update education : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, dto.MessageResponseWithID{
			Message:   "Education updated successfully",
			ProfileID: value,
		})
	}
}
