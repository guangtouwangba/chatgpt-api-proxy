package httphelper

import (
	"chatgpt-api-proxy/internal/constant"
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

func WrapperError(c *gin.Context, err constant.BaseError) {
	WrapperResponse(c, &BaseResponse{
		Code:    constant.ErrorCodesToHTTPStatusCodes[err.Code],
		Message: err.Message,
		Data:    nil,
	})
}
