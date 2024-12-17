package filesystem

import (
	"github.com/gogo/protobuf/proto"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

func init() {
	protoutils.RegisterDecoderFunc("type.googleapis.com/filesystem.ListResp", func(b []byte) (proto.Message, error) {
		resp := new(ListResp)
		err := proto.Unmarshal(b, resp)
		if err != nil {
			return nil, err
		}

		return resp, nil
	})
}

func (e *Entry) IsDir() bool {
	return e.GetType() == EntryTypeDir
}
