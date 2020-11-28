package routes

import (
	"go-web-console/exceptions"

	"github.com/gin-gonic/gin"
)

func GetGin() *gin.Engine {
	r := gin.New()
	r.Use(exceptions.HandleErrorMiddleware)
	r = initV1ApiGroup(r)
	r = initV2ApiGroup(r)
	return r
}

func initV1ApiGroup(r *gin.Engine) *gin.Engine {
	v1ws := r.Group("/ws/v1")
	v1ws.GET("/ping", v1WsPing)
	v1ws.GET("/links", v1wsLink)
	v1api := r.Group("/api/v1")
	v1api.GET("/ping", v1ApiPing)
	v1api.POST("/links", v1ApiCreateLink)
	v1api.GET("/links", v1ApiGetLinks)
	return r
}

func initV2ApiGroup(r *gin.Engine) *gin.Engine {
	return r
}
