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

// CreateCertificateHandler handles HTTP requests to add certificates details to a user profile.
func CreateCertificateHandler(ctx context.Context, certificateSvc service.Service) func(http.ResponseWriter, *http.Request) {
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

		req, err := decodeCreateCertificateRequest(r)
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

		profileID, err = certificateSvc.CreateCertificate(ctx, req, profileID, userID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to create certificate : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponseWithID{
			Message:   "Certificate(s) added successfully",
			ProfileID: profileID,
		})
	}
}

// ListCertificatesHandler handles HTTP requests to get certificates details of a user profile.
func ListCertificatesHandler(ctx context.Context, certificateSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParamsByID(r, constants.ProfileID)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		filter, err := helpers.DecodeCertificateRequest(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadRequest, errors.ErrDecodeRequest)
			zap.S().Error(err)
			return
		}

		cetificateResp, err := certificateSvc.ListCertificates(ctx, profileID, filter)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailespecsFetch)
			zap.S().Error("Unable to fetch certificate : ", err, "for profile id : ", profileID)
			return
		}

		if len(cetificateResp) == 0 {
			cetificateResp = []specs.CertificateResponse{}
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.ResponseCertificate{
			Certificates: cetificateResp,
		})
	}
}

// UpdateCertificateHandler returns an HTTP handler that updates certificates using profileSvc.
func UpdateCertificateHandler(ctx context.Context, certificateSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, certID, err := helpers.GetMultipleParams(r)
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

		req, err := decodeUpdateCertificateRequest(r)
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

		updatedResp, err := certificateSvc.UpdateCertificate(ctx, profileID, certID, userID, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to update certificate : ", err, "for profile id : ", profileID, "certificate id : ", certID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponseWithID{
			Message:   "Certificate updated successfully",
			ProfileID: updatedResp,
		})
	}
}

// DeleteCertificatesHandler returns an HTTP handler that deletes certificates using profileSvc.
func DeleteCertificatesHandler(ctx context.Context, certificateSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, certificateID, err := helpers.GetMultipleParams(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Error getting profile id and certificate id : ", err)
			return
		}
		// call the service
		err = certificateSvc.DeleteCertificate(ctx, profileID, certificateID)
		if err != nil {
			if err == errors.ErrNoData {
				middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
					Message: constants.NoResourceFound,
				})
				return
			}
			middleware.ErrorResponse(w, http.StatusBadGateway, errors.ErrFailedToDelete)
			zap.S().Error("Unable to delete certificate : ", err, "for profile id : ", profileID, "certificate id : ", certificateID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK, specs.MessageResponse{
			Message: "Certificate deleted successfully",
		})
	}
}
