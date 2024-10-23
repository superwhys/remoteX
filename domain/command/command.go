package command

import (
	"encoding/json"

	"github.com/superwhys/remoteX/pkg/protoutils"
)

func StrArg(val string) Command_Arg {
	return Command_Arg{
		Value: &Command_Arg_StrValue{StrValue: val},
	}
}

func BoolArg(val bool) Command_Arg {
	return Command_Arg{
		Value: &Command_Arg_BoolValue{BoolValue: val},
	}
}

func IntArg(val int64) Command_Arg {
	return Command_Arg{
		Value: &Command_Arg_IntValue{IntValue: val},
	}
}

func EmptyCommand() *Command {
	return &Command{
		Type: Empty,
		Args: map[string]Command_Arg{
			"type": {Value: &Command_Arg_StrValue{StrValue: "empty"}},
		},
	}
}

func (m *Ret) MarshalJSON() ([]byte, error) {
	respValue := struct {
		Resp   any    `json:"resp,omitempty"`
		ErrNo  uint64 `json:"errNo,omitempty"`
		ErrMsg string `json:"errMsg,omitempty"`
	}{
		ErrNo:  m.ErrNo,
		ErrMsg: m.ErrMsg,
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
