package api

import (
	"github.com/gin-gonic/gin"
	"github.com/go-puzzles/puzzles/pgin"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/pkg/errorutils"
	"github.com/superwhys/remoteX/server"
)

type ApiService[R any] interface {
	Handle(c *gin.Context, srv *server.RemoteXServer) (resp R, err error)
}

type Ret struct {
	ErrNo  uint64 `json:"err_no"`
	ErrMsg string `json:"err_msg"`
	Data   any    `json:"data"`
}

type RemoteXAPI struct {
	srv *server.RemoteXServer
}

func NewRemoteXAPI(srv *server.RemoteXServer) *RemoteXAPI {
	return &RemoteXAPI{
		srv: srv,
	}
}

func MountApiService[AS ApiService[R], R any](srv *server.RemoteXServer) gin.HandlerFunc {
	return pgin.RequestResponseHandler(func(c *gin.Context, req *AS) (*Ret, error) {
		resp, err := (*req).Handle(c, srv)

		ret := &Ret{Data: resp}
		if err != nil {
			var be *errorutils.RemoteXError
			if errors.As(err, &be) {
				ret.ErrNo = be.Code()
			} else {
				ret.ErrNo = 0
			}

			ret.ErrMsg = err.Error()
		}

		return ret, err
	})
}

func (a *RemoteXAPI) SetupHttpServer() *gin.Engine {
	router := pgin.Default()

	router.GET("/node", MountApiService[getAllNodes](a.srv))
	router.GET("/list/dir", MountApiService[listDir](a.srv))
	router.POST("/sync/pull", MountApiService[syncPull](a.srv))
	router.POST("/sync/push", MountApiService[syncPush](a.srv))
	router.POST("/tunnel/forward", MountApiService[tunnelForward](a.srv))
	router.GET("/tunnel/list", MountApiService[listTunnel](a.srv))
	router.POST("/tunnel/close", MountApiService[closeTunnel](a.srv))

	// nodeRouter is used to handle remote commands
	nodeRouter := router.Group("/node/:nodeId")
	{
		nodeRouter.GET("", MountApiService[getNode](a.srv))
		nodeRouter.GET("/list/dir", MountApiService[listRemoteDir](a.srv))
	}

	return router
}
