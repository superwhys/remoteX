package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/pgin"
	"github.com/superwhys/remoteX/server"
)

type RemoteXAPI struct {
	srv *server.RemoteXServer
}

func NewRemoteXAPI(srv *server.RemoteXServer) *RemoteXAPI {
	return &RemoteXAPI{
		srv: srv,
	}
}

func (a *RemoteXAPI) SetupHttpServer() *gin.Engine {
	router := pgin.Default()

	router.GET("/node", a.getAllNodes())
	router.GET("/list/dir", pgin.RequestHandler(a.listDir))
	router.POST("/sync/pull", pgin.RequestHandler(a.pullEntry))
	router.POST("/sync/push", pgin.RequestHandler(a.pushEntry))

	nodeRouter := router.Group("/node/:nodeId")
	{
		nodeRouter.GET("", pgin.RequestHandler(a.getNode))
		nodeRouter.GET("/list/dir", pgin.RequestHandler(a.listRemoteDir))
	}

	return router
}
