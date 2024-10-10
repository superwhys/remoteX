package server

import (
	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/pgin"
)

func (s *RemoteXServer) SetupHttpServer() *gin.Engine {
	engine := pgin.Default()

	nodeRouter := engine.Group("node")
	{
		nodeRouter.GET("", s.getAllNodes())
		nodeRouter.GET("/:nodeId", pgin.RequestHandler(s.getNode))
		nodeRouter.GET("/list/dir", pgin.RequestHandler(s.listDir))
		nodeRouter.GET("/:nodeId/list/dir", pgin.RequestHandler(s.listRemoteDir))
	}

	return engine
}
