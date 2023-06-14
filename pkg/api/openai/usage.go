package api

import (
	"chatgpt-api-proxy/internal/db"
	"chatgpt-api-proxy/internal/db/repository"
	"chatgpt-api-proxy/pkg/httphelper"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitUsageRouter(r *gin.Engine) {
	api := r.Group("/api/openai")
	api.GET("usage/:id", HandleGetUsage)
}

func HandleGetUsage(c *gin.Context) {
	id := c.Param("id")
	database := db.GetDB()
	repo := repository.NewGormOpenAIUsageRepository(database)
	usage, err := repo.FindOne(id, "")
	if err != nil {
		httphelper.WrapperError(c, httphelper.NewAPIError(http.StatusInternalServerError, err.Error()))
		return
	}

	httphelper.WrapperSuccess(c, usage)
}
