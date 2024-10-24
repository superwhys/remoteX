package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/pgin"
)

func (s *RemoteXServer) SetupHttpServer() *gin.Engine {
	router := pgin.Default()

	router.GET("/node", s.getAllNodes())
	router.GET("/list/dir", pgin.RequestHandler(s.listDir))
	router.POST("/sync/pull", pgin.RequestHandler(s.pullEntry))
	router.POST("/sync/push", pgin.RequestHandler(s.pushEntry))

	nodeRouter := router.Group("/node/:nodeId")
	{
		nodeRouter.GET("", pgin.RequestHandler(s.getNode))
		nodeRouter.GET("/list/dir", pgin.RequestHandler(s.listRemoteDir))
	}

	return router
}
