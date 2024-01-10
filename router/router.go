package router

import (
	"BBoltViewer/cmd"
	"BBoltViewer/controller"
	"BBoltViewer/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func Init(c *cmd.Command) *gin.Engine {
	r := gin.New()
	if !c.Debug() {
		gin.SetMode(gin.ReleaseMode)
	}
	r.Use(
		middleware.RateLimit(time.Second, 100, 10),
		middleware.Cors(),
		gin.Logger(),
		gin.Recovery(),
	)
	if !c.Http {
		r.Use(middleware.LoadTls(c.ServerPort))
	}
	r.POST("/connect", controller.Connect())
	r.POST("/close", controller.Close())
	r.POST("/data_from_bucket", controller.DataFromBucket())
	r.POST("/set", controller.Set())
	r.POST("/get", controller.Get())
	r.POST("/delete", controller.Delete())
	return r
}
