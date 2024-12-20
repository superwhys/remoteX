package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/superwhys/remoteX/internal/proto/ext"

	"github.com/gogo/protobuf/gogoproto"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/protoc-gen-gogo/descriptor"
	"github.com/gogo/protobuf/vanity"
	"github.com/gogo/protobuf/vanity/command"
)

func main() {
	req := command.Read()
	files := req.GetProtoFile()
	files = vanity.FilterFiles(files, vanity.NotGoogleProtobufDescriptorProto)

	vanity.ForEachFile(files, TurnOnProtoSizerAll)
	vanity.ForEachFile(files, vanity.TurnOffGoEnumPrefixAll)
	vanity.ForEachFile(files, vanity.TurnOffGoUnrecognizedAll)
	vanity.ForEachFile(files, vanity.TurnOffGoUnkeyedAll)
	vanity.ForEachFile(files, vanity.TurnOffGoSizecacheAll)
	vanity.ForEachFile(files, vanity.TurnOnMarshalerAll)
	vanity.ForEachFile(files, vanity.TurnOnUnmarshalerAll)
	vanity.ForEachEnumInFiles(files, HandleCustomEnumExtensions)
	vanity.ForEachFile(files, SetPackagePrefix("github.com/superwhys/remoteX"))
	vanity.ForEachFile(files, HandleFile)
	vanity.ForEachFieldInFilesExcludingExtensions(files, TurnOffNullableForMessages)

	resp := command.Generate(req)
	command.Write(resp)
}

func TurnOnProtoSizerAll(file *descriptor.FileDescriptorProto) {
	vanity.SetBoolFileOption(gogoproto.E_ProtosizerAll, true)(file)
}

func TurnOffNullableForMessages(field *descriptor.FieldDescriptorProto) {
	if !vanity.FieldHasBoolExtension(field, gogoproto.E_Nullable) {
		_, hasCustomType := GetFieldStringExtension(field, gogoproto.E_Customtype)
		if field.IsMessage() || hasCustomType {
			vanity.SetBoolFieldOption(gogoproto.E_Nullable, false)(field)
		}
	}
}

func HandleCustomEnumExtensions(enum *descriptor.EnumDescriptorProto) {
	for _, field := range enum.Value {
		if field == nil {
			continue
		}
		if field.Options == nil {
			field.Options = &descriptor.EnumValueOptions{}
		}
		customName := gogoproto.GetEnumValueCustomName(field)
		if customName != "" {
			continue
		}
		if v, ok := GetEnumValueStringExtension(field, ext.E_Enumgoname); ok {
			SetEnumValueStringFieldOption(field, gogoproto.E_EnumvalueCustomname, v)
		} else {
			SetEnumValueStringFieldOption(field, gogoproto.E_EnumvalueCustomname, toCamelCase(*field.Name, true))
		}
	}
}

func SetPackagePrefix(prefix string) func(file *descriptor.FileDescriptorProto) {
	return func(file *descriptor.FileDescriptorProto) {
		if file.Options.GoPackage == nil {
			pkg, _ := filepath.Split(file.GetName())
			fullPkg := prefix + "/" + strings.TrimSuffix(pkg, "/")
			file.Options.GoPackage = &fullPkg
		}
	}
}

func toCamelCase(input string, firstUpper bool) string {
	runes := []rune(strings.ToLower(input))
	outputRunes := make([]rune, 0, len(runes))

	nextUpper := false
	for i, rune := range runes {
		if rune == '_' {
			nextUpper = true
			continue
		}
		if (firstUpper && i == 0) || nextUpper {
			rune = unicode.ToUpper(rune)
			nextUpper = false
		}
		outputRunes = append(outputRunes, rune)
	}
	return string(outputRunes)
}

func SetStringFieldOption(field *descriptor.FieldDescriptorProto, extension *proto.ExtensionDesc, value string) {
	if _, ok := GetFieldStringExtension(field, extension); ok {
		return
	}
	if field.Options == nil {
		field.Options = &descriptor.FieldOptions{}
	}
	if err := proto.SetExtension(field.Options, extension, &value); err != nil {
		panic(err)
	}
}

func SetEnumValueStringFieldOption(field *descriptor.EnumValueDescriptorProto, extension *proto.ExtensionDesc, value string) {
	if _, ok := GetEnumValueStringExtension(field, extension); ok {
		return
	}
	if field.Options == nil {
		field.Options = &descriptor.EnumValueOptions{}
	}
	if err := proto.SetExtension(field.Options, extension, &value); err != nil {
		panic(err)
	}
}

func GetEnumValueStringExtension(enumValue *descriptor.EnumValueDescriptorProto, extension *proto.ExtensionDesc) (string, bool) {
	if enumValue.Options == nil {
		return "", false
	}
	value, err := proto.GetExtension(enumValue.Options, extension)
	if err != nil {
		return "", false
	}
	if value == nil {
		return "", false
	}
	if v, ok := value.(*string); !ok || v == nil {
		return "", false
	} else {
		return *v, true
	}
}

func GetFieldStringExtension(field *descriptor.FieldDescriptorProto, extension *proto.ExtensionDesc) (string, bool) {
	if field.Options == nil {
		return "", false
	}
	value, err := proto.GetExtension(field.Options, extension)
	if err != nil {
		return "", false
	}
	if value == nil {
		return "", false
	}
	if v, ok := value.(*string); !ok || v == nil {
		return "", false
	} else {
		return *v, true
	}
}

func GetFieldBooleanExtension(field *descriptor.FieldDescriptorProto, extension *proto.ExtensionDesc) (bool, bool) {
	if field.Options == nil {
		return false, false
	}
	value, err := proto.GetExtension(field.Options, extension)
	if err != nil {
		return false, false
	}
	if value == nil {
		return false, false
	}
	if v, ok := value.(*bool); !ok || v == nil {
		return false, false
	} else {
		return *v, true
	}
}

func GetMessageBoolExtension(msg *descriptor.DescriptorProto, extension *proto.ExtensionDesc) (bool, bool) {
	if msg.Options == nil {
		return false, false
	}
	value, err := proto.GetExtension(msg.Options, extension)
	if err != nil {
		return false, false
	}
	if value == nil {
		return false, false
	}
	val, ok := value.(*bool)
	if !ok || val == nil {
		return false, false
	}
	return *val, true
}

func HandleFile(file *descriptor.FileDescriptorProto) {
	vanity.ForEachMessageInFiles([]*descriptor.FileDescriptorProto{file}, HandleCustomExtensions(file))
}

func HandleCustomExtensions(file *descriptor.FileDescriptorProto) func(msg *descriptor.DescriptorProto) {
	return func(msg *descriptor.DescriptorProto) {
		vanity.ForEachField([]*descriptor.DescriptorProto{msg}, func(field *descriptor.FieldDescriptorProto) {
			if field.Options == nil {
				field.Options = &descriptor.FieldOptions{}
			}
			deprecated := field.Options.Deprecated != nil && *field.Options.Deprecated == true

			if field.Type != nil && *field.Type == descriptor.FieldDescriptorProto_TYPE_INT32 {
				SetStringFieldOption(field, gogoproto.E_Casttype, "int")
			}

			if field.TypeName != nil && *field.TypeName == ".google.protobuf.Timestamp" {
				vanity.SetBoolFieldOption(gogoproto.E_Stdtime, true)(field)
			}

			if _, ok := GetFieldBooleanExtension(field, ext.E_Nullable); ok {
				vanity.SetBoolFieldOption(gogoproto.E_Nullable, true)(field)
			}

			if goName, ok := GetFieldStringExtension(field, ext.E_Goname); ok {
				SetStringFieldOption(field, gogoproto.E_Customname, goName)
			} else if deprecated {
				SetStringFieldOption(field, gogoproto.E_Customname, "Deprecated"+toCamelCase(*field.Name, true))
			}

			if goType, ok := GetFieldStringExtension(field, ext.E_Gotype); ok {
				SetStringFieldOption(field, gogoproto.E_Customtype, goType)
			}

			if val, ok := GetFieldBooleanExtension(field, ext.E_NodeId); ok && val {
				if *file.Options.GoPackage != "github.com/superwhys/remoteX/pkg/common" {
					SetStringFieldOption(field, gogoproto.E_Customtype, "github.com/superwhys/remoteX/pkg/common.NodeID")
				} else {
					SetStringFieldOption(field, gogoproto.E_Customtype, "NodeID")
				}
			}

			if yamlValue, ok := GetFieldStringExtension(field, ext.E_Yaml); ok {
				SetStringFieldOption(field, gogoproto.E_Moretags, fmt.Sprintf(`yaml:"%s"`, yamlValue))
			} else {
				SetStringFieldOption(field, gogoproto.E_Moretags, `yaml:"-"`)
			}

			if jsonValue, ok := GetFieldStringExtension(field, ext.E_Json); ok {
				SetStringFieldOption(field, gogoproto.E_Jsontag, jsonValue)
			} else if deprecated {
				SetStringFieldOption(field, gogoproto.E_Jsontag, "-")
			} else {
				SetStringFieldOption(field, gogoproto.E_Jsontag, toCamelCase(*field.Name, false))
			}

			current := ""
			if v, ok := GetFieldStringExtension(field, gogoproto.E_Moretags); ok {
				current = v
			}

			SetStringFieldOption(field, gogoproto.E_Moretags, current)
		})
	}
}
