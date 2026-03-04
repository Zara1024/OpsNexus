package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/response"
)

func handleAppError(c *gin.Context, err error) {
	var appErr *apperrors.AppError
	if errors.As(err, &appErr) {
		response.Error(c, appErr.HTTPCode, appErr.Code, appErr.Message)
		return
	}
	response.Error(c, http.StatusInternalServerError, apperrors.ErrInternal.Code, apperrors.ErrInternal.Message)
}
