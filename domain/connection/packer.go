package connection

import (
	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/pkg/counter"
	"github.com/superwhys/remoteX/pkg/limiter"
	"github.com/superwhys/remoteX/pkg/protoutils"
	"github.com/superwhys/remoteX/pkg/tracker"
)

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

var _ Stream = (*TrackerStream)(nil)

type TrackerStream struct {
	Stream
	rd *tracker.TrackerReader
	wr *tracker.TrackerWriter
	protoutils.ProtoMessageReader
	protoutils.ProtoMessageWriter
}

func PackTrackerStream(manager *tracker.Manager, stream Stream) *TrackerStream {
	rd := tracker.NewTrackerReader(stream, manager)
	wr := tracker.NewTrackerWriter(stream, manager)

	return &TrackerStream{
		Stream:             stream,
		rd:                 rd,
		wr:                 wr,
		ProtoMessageReader: protoutils.NewProtoReader(rd),
		ProtoMessageWriter: protoutils.NewProtoWriter(wr),
	}
}

func (t *TrackerStream) Pack(stream Stream) Stream {
	return PackTrackerStream(tracker.Trackermanager, stream)
}

func (t *TrackerStream) Read(p []byte) (n int, err error) {
	return t.rd.Read(p)
}

func (t *TrackerStream) Write(p []byte) (n int, err error) {
	return t.wr.Write(p)
}

func (t *TrackerStream) ReadMessage(message proto.Message) error {
	return t.ProtoMessageReader.ReadMessage(message)
}

func (t *TrackerStream) WriteMessage(m proto.Message) error {
	return t.ProtoMessageWriter.WriteMessage(m)
}

func (t *TrackerStream) Close() (err error) {
	return t.Stream.Close()
}

var _ Stream = (*LimiterStream)(nil)

type LimiterStream struct {
	Stream
	rd *limiter.LimitReader
	wr *limiter.LimitWriter
	protoutils.ProtoMessageReader
	protoutils.ProtoMessageWriter
}

func PackLimiterStream(stream Stream, limiter *limiter.Limiter) *LimiterStream {
	rd, wr := limiter.GetNodeRateLimiter(stream)
	return &LimiterStream{
		Stream:             stream,
		rd:                 rd,
		wr:                 wr,
		ProtoMessageReader: protoutils.NewProtoReader(rd),
		ProtoMessageWriter: protoutils.NewProtoWriter(wr),
	}
}

func (l *LimiterStream) Pack(stream Stream) Stream {
	return PackLimiterStream(stream, limiter.StreamLimiter)
}

// Read rewrite the method to use LimiterReader
func (l *LimiterStream) Read(p []byte) (n int, err error) {
	return l.rd.Read(p)
}

// Write rewrite the method to use LimiterWriter
func (l *LimiterStream) Write(p []byte) (n int, err error) {
	return l.wr.Write(p)
}

func (l *LimiterStream) ReadMessage(message proto.Message) error {
	return l.ProtoMessageReader.ReadMessage(message)
}

func (l *LimiterStream) WriteMessage(m proto.Message) error {
	return l.ProtoMessageWriter.WriteMessage(m)
}

func (l *LimiterStream) Close() (err error) {
	return l.Stream.Close()
}

var _ Stream = (*CounterStream)(nil)

type CounterStream struct {
	Stream
	protoutils.ProtoMessageReader
	protoutils.ProtoMessageWriter

	rd *counter.CountingReader
	wr *counter.CountingWriter
}

func PackCounterStream(stream Stream) *CounterStream {
	crd := &counter.CountingReader{Reader: stream}
	cwr := &counter.CountingWriter{Writer: stream}

	return &CounterStream{
		Stream:             stream,
		rd:                 crd,
		wr:                 cwr,
		ProtoMessageReader: protoutils.NewProtoReader(crd),
		ProtoMessageWriter: protoutils.NewProtoWriter(cwr),
	}
}

func (cc *CounterStream) Pack(stream Stream) Stream {
	return PackCounterStream(stream)
}

// Read rewrite the method to use CountingReader
func (cc *CounterStream) Read(p []byte) (n int, err error) {
	return cc.rd.Read(p)
}

// Write rewrite the method to use CountingWriter
func (cc *CounterStream) Write(p []byte) (n int, err error) {
	return cc.wr.Write(p)
}

func (cc *CounterStream) ReadMessage(message proto.Message) error {
	return cc.ProtoMessageReader.ReadMessage(message)
}

func (cc *CounterStream) WriteMessage(m proto.Message) error {
	return cc.ProtoMessageWriter.WriteMessage(m)
}

func (cc *CounterStream) Close() error {
	return cc.Stream.Close()
}
