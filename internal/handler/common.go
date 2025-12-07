package handler

import (
	v1 "go-nunu/api/v1"
	"go-nunu/internal/service"
	"go-nunu/pkg/aws"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CommonHandler struct {
	*Handler
	CommonService service.CommonService
}

func NewCommonHandler(handler *Handler, commonService service.CommonService, r2Client *aws.CloudflareR2) *CommonHandler {
	return &CommonHandler{
		Handler:       handler,
		CommonService: commonService,
	}
}

func (h *CommonHandler) UploadPresignedUrl(ctx *gin.Context) {
	var req v1.UploadPresignedUrlRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	preSignedUrl, endpointUrl, err := h.CommonService.UploadPresignedUrl(req.FileExt, req.UploadScene)

	if err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	v1.HandleSuccess(ctx, v1.UploadPresignedUrlResponseData{
		PreSignedUrl: preSignedUrl,
		EndpointUrl:  endpointUrl,
	})
}
