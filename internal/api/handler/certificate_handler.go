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

// CreateCertificateHandler handles HTTP requests to add certificates details to a user profile.
func CreateCertificateHandler(ctx context.Context, profileSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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

		profileID, err := profileSvc.CreateCertificate(ctx, req)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to create certificate : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusCreated, dto.MessageResponseWithID{
			Message:   "Certificate(s) added successfully",
			ProfileID: profileID,
		})
	}
}

// GetCertificatesHandler handles HTTP requests to get certificates details of a user profile.
func ListCertificatesHandler(ctx context.Context, certificateSvc service.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		profileID, err := helpers.GetParams(r)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error(err)
			return
		}

		id, err := helpers.ConvertStringToInt(profileID)
		if err != nil {
			zap.S().Error("error to get education : ", err, " for profile id : ", profileID)
			return
		}

		values, err := certificateSvc.ListCertificates(ctx, id)
		if err != nil {
			middleware.ErrorResponse(w, http.StatusBadGateway, err)
			zap.S().Error("Unable to fetch certificate : ", err, "for profile id : ", profileID)
			return
		}

		middleware.SuccessResponse(w, http.StatusOK,
			dto.ResponseCertificate{
				Certificates: values,
			})
	}
}
