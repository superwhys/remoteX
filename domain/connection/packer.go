package connection

type packer interface {
	Pack(stream Stream) Stream
}

var (
	packers = []packer{
		&TrackerStream{},
		&LimiterStream{},
		&CounterStream{},
	}
)

func PackStream(stream Stream) Stream {
	for _, p := range packers {
		stream = p.Pack(stream)
	}
	return stream
}
