package command

import (
	"encoding/json"

	"github.com/superwhys/remoteX/pkg/protoutils"
)

func EmptyCommand() *Command {
	return &Command{
		Type: Empty,
		Args: map[string]string{"type": "empty"},
	}
}

func (m *Ret) MarshalJSON() ([]byte, error) {
	respValue := struct {
		Command *Command `json:"command,omitempty"`
		Resp    any      `json:"resp,omitempty"`
		ErrNo   uint64   `json:"errNo,omitempty"`
		ErrMsg  string   `json:"errMsg,omitempty"`
	}{
		Command: m.Command,
		ErrNo:   m.ErrNo,
		ErrMsg:  m.ErrMsg,
	}

	if m.Resp != nil {
		pm, err := protoutils.DecodeAnyProto(m.Resp)
		if err != nil {
			respValue.Resp = err.Error()
		} else {
			respValue.Resp = pm
		}
	}

	return json.Marshal(respValue)
}
