package command

import (
	"encoding/json"
	"fmt"
	
	"github.com/superwhys/remoteX/pkg/protoutils"
)

func EmptyCommand() *Command {
	return &Command{
		Type: Empty,
		Args: map[string]string{"type": "empty"},
	}
}

func (m *Ret) MarshalJSON() ([]byte, error) {
	respValue := map[string]any{
		"command": m.Command,
		"errNo":   m.ErrNo,
		"errMsg":  m.ErrMsg,
	}
	
	if m.Resp != nil {
		fmt.Println(m.Resp.TypeUrl)
		pm, err := protoutils.DecodeAnyProto(m.Resp)
		if err != nil {
			respValue["resp"] = err.Error()
		} else {
			respValue["resp"] = pm
		}
	}
	
	return json.Marshal(respValue)
}
