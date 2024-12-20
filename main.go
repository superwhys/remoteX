package main

import (
	"context"

	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/api"
	"github.com/superwhys/remoteX/config"
	"github.com/superwhys/remoteX/server"
	"github.com/thejerf/suture/v4"

	httppuzzle "github.com/go-puzzles/puzzles/cores/puzzles/http-puzzle"
	pprofpuzzle "github.com/go-puzzles/puzzles/cores/puzzles/pprof-puzzle"
)

var (
	port       = pflags.Int("port", 0, "server run port")
	configFlag = pflags.Struct("conf", (*config.Config)(nil), "remoteX config")
)

func main() {
	pflags.SetStructParseTagName("yaml")
	pflags.Parse()

	conf := new(config.Config)
	plog.PanicError(configFlag(conf))

	opt, err := server.InitOption(conf)
	plog.PanicError(err)

	remoteXServer := server.NewRemoteXServer(opt)
	remoteXApi := api.NewRemoteXAPI(remoteXServer)

	core := cores.NewPuzzleCore(
		pprofpuzzle.WithCorePprof(),
		httppuzzle.WithCoreHttpPuzzle("/api", remoteXApi.SetupHttpServer()),
		cores.WithDaemonNameWorker("RemoteX", func(ctx context.Context) error {
			mainService := suture.NewSimple("main")
			mainService.Add(remoteXServer)
			return mainService.Serve(ctx)
		}),
	)

	plog.PanicError(cores.Start(core, port()))
}
