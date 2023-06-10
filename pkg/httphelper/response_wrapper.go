package httphelper

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func WrapperResponse(c *gin.Context, response *BaseResponse) {
	c.JSON(response.Code, response)
}

func WrapperSuccess(c *gin.Context, data interface{}) {
	WrapperResponse(c, &BaseResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	})
}

func WrapperError(c *gin.Context, err error) {
	apiError := &APIError{}
	if errors.As(err, &apiError) {
		WrapperResponse(c, &BaseResponse{
			Code:    apiError.Code,
			Message: apiError.Message,
			Data:    nil,
		})
		return
	}
	WrapperResponse(c, &BaseResponse{
		Code:    http.StatusInternalServerError,
		Message: err.Error(),
		Data:    nil,
	})
}
