package protoutils

type ProtoMessageReadWriter interface {
	ProtoMessageReader
	ProtoMessageWriter
}
