package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// First generate extensions using standard proto compiler.
//go:generate protoc -I .. -I . -I $GOPATH/src/github.com/gogo/protobuf/protobuf --gogofast_out=Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/protoc-gen-gogo/descriptor,paths=source_relative:ext ext.proto

// Then build our vanity compiler that uses the new extensions
//go:generate go build -o scripts/protoc-gen-remotex scripts/protoc_plugin.go

//go:generate go run generate.go config pkg/protocol domain/node domain/connection domain/command internal/fs

func main() {
	goPath := os.Getenv("GOPATH")
	for _, path := range os.Args[1:] {
		matches, err := filepath.Glob(filepath.Join(path, "*proto"))
		if err != nil {
			log.Fatal(err)
		}
		log.Println(path, "returned:", matches)

		args := []string{
			"-I", "..",
			"-I", ".",
			"-I", fmt.Sprintf("%v/src/github.com/gogo/protobuf/protobuf", goPath),
			"--plugin=protoc-gen-goremotex=scripts/protoc-gen-remotex",
			"--goremotex_out", "Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types,paths=source_relative:../..",
		}
		args = append(args, matches...)
		cmd := exec.Command("protoc", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal("Failed generating", path)
		}
	}
}
