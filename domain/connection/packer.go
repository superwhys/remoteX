package connection

import (
	"github.com/superwhys/remoteX/pkg/limiter"
	"github.com/superwhys/remoteX/pkg/tracker"
)

type PackOpts struct {
	Limiter        *limiter.Limiter
	TrackerManager *tracker.Manager
}

type packer interface {
	Pack(stream Stream, opts *PackOpts) Stream
}

var (
	packers = []packer{
		&TrackerStream{},
		&LimiterStream{},
		&CounterStream{},
	}
)

func PackStream(stream Stream, opts *PackOpts) Stream {
	for _, p := range packers {
		stream = p.Pack(stream, opts)
	}
	return stream
}
