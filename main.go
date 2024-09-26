package main

import (
	"context"

	"github.com/go-puzzles/puzzles/cores"
	"github.com/go-puzzles/puzzles/pflags"
	"github.com/go-puzzles/puzzles/plog"
	"github.com/superwhys/remoteX/config"
	"github.com/superwhys/remoteX/server"
	"github.com/thejerf/suture/v4"
)

var (
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

	core := cores.NewPuzzleCore(
		cores.WithDaemonNameWorker("RemoteX", func(ctx context.Context) error {
			mainService := suture.NewSimple("main")
			mainService.Add(remoteXServer)
			return mainService.Serve(ctx)
		}),
	)

	plog.PanicError(cores.Run(core))
}
