// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package template

import (
	"net/http"

	"github.com/LF-Engineering/lfx-kit/auth"
	"github.com/communitybridge/easycla/cla-backend-go/events"
	v1Events "github.com/communitybridge/easycla/cla-backend-go/events"
	v1Models "github.com/communitybridge/easycla/cla-backend-go/gen/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/models"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations"
	"github.com/communitybridge/easycla/cla-backend-go/gen/v2/restapi/operations/template"
	log "github.com/communitybridge/easycla/cla-backend-go/logging"
	v1Template "github.com/communitybridge/easycla/cla-backend-go/template"
	"github.com/communitybridge/easycla/cla-backend-go/utils"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/jinzhu/copier"
)

// Configure API call
func Configure(api *operations.EasyclaAPI, service v1Template.Service, eventsService v1Events.Service) {
	// Retrieve a list of available templates
	api.TemplateGetTemplatesHandler = template.GetTemplatesHandlerFunc(func(params template.GetTemplatesParams, user *auth.User) middleware.Responder {

		templates, err := service.GetTemplates(params.HTTPRequest.Context())
		if err != nil {
			return template.NewGetTemplatesBadRequest().WithPayload(errorResponse(err))
		}
		response := []models.Template{}
		err = copier.Copy(&response, templates)
		if err != nil {
			return template.NewGetTemplatesInternalServerError().WithPayload(errorResponse(err))
		}
		return template.NewGetTemplatesOK().WithPayload(response)
	})

	api.TemplateCreateCLAGroupTemplateHandler = template.CreateCLAGroupTemplateHandlerFunc(func(params template.CreateCLAGroupTemplateParams, user *auth.User) middleware.Responder {
		utils.SetAuthUserProperties(user, params.XUSERNAME, params.XEMAIL)
		input := &v1Models.CreateClaGroupTemplate{}
		err := copier.Copy(input, &params.Body)
		if err != nil {
			return template.NewGetTemplatesInternalServerError().WithPayload(errorResponse(err))
		}
		pdfUrls, err := service.CreateCLAGroupTemplate(params.HTTPRequest.Context(), params.ClaGroupID, input)
		if err != nil {
			log.Warnf("Error generating PDFs from provided templates, error: %v", err)
			return template.NewGetTemplatesBadRequest().WithPayload(errorResponse(err))
		}
		eventsService.LogEvent(&events.LogEventArgs{
			EventType:  events.CLATemplateCreated,
			ProjectID:  params.ClaGroupID,
			LfUsername: user.UserName,
			EventData:  &events.CLATemplateCreatedEventData{},
		})

		response := &models.TemplatePdfs{}
		err = copier.Copy(response, pdfUrls)
		if err != nil {
			return template.NewGetTemplatesInternalServerError().WithPayload(errorResponse(err))
		}
		return template.NewCreateCLAGroupTemplateOK().WithPayload(response)
	})

	api.TemplateTemplatePreviewHandler = template.TemplatePreviewHandlerFunc(func(params template.TemplatePreviewParams, user *auth.User) middleware.Responder {
		var param v1Models.CreateClaGroupTemplate
		err := copier.Copy(&param, &params.TemplatePreviewInput)
		if err != nil {
			return writeResponse(http.StatusInternalServerError, runtime.JSONMime, runtime.JSONProducer(), errorResponse(err))
		}
		pdf, err := service.CreateTemplatePreview(&param, params.TemplateFor)
		if err != nil {
			log.Warnf("Error generating PDFs from provided templates, error: %v", err)
			return writeResponse(http.StatusBadRequest, runtime.JSONMime, runtime.JSONProducer(), errorResponse(err))
		}
		return middleware.ResponderFunc(func(rw http.ResponseWriter, pr runtime.Producer) {
			rw.WriteHeader(http.StatusOK)
			_, err := rw.Write(pdf)
			if err != nil {
				log.Warnf("Error writing pdf, error: %v", err)
			}
		})
	})

	api.TemplateGetCLATemplatePreviewHandler = template.GetCLATemplatePreviewHandlerFunc(func(params template.GetCLATemplatePreviewParams) middleware.Responder {
		pdf, err := service.GetCLATemplatePreview(params.HTTPRequest.Context(), params.ClaGroupID, params.ClaType, *params.Watermark)
		if err != nil {
			log.Warnf("Error getting PDFs for provided clagroupID : %s, error: %v", params.ClaGroupID, err)
			return writeResponse(http.StatusBadRequest, runtime.JSONMime, runtime.JSONProducer(), errorResponse(err))
		}
		return middleware.ResponderFunc(func(rw http.ResponseWriter, pr runtime.Producer) {
			rw.WriteHeader(http.StatusOK)
			_, err := rw.Write(pdf)
			if err != nil {
				log.Warnf("Error writing pdf, error: %v", err)
			}
		})
	})
}

type codedResponse interface {
	Code() string
}

func errorResponse(err error) *models.ErrorResponse {
	code := ""
	if e, ok := err.(codedResponse); ok {
		code = e.Code()
	}

	e := models.ErrorResponse{
		Code:    code,
		Message: err.Error(),
	}

	return &e
}

func writeResponse(httpStatus int, contentType string, contentProducer runtime.Producer, data interface{}) middleware.Responder {
	return middleware.ResponderFunc(func(rw http.ResponseWriter, pr runtime.Producer) {
		rw.Header().Set(runtime.HeaderContentType, contentType)
		rw.WriteHeader(httpStatus)
		err := contentProducer.Produce(rw, data)
		if err != nil {
			log.Warnf("failed to write data. error = %v", err)
		}
	})
}
