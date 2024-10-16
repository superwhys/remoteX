package filesync

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/stretchr/testify/assert"
)

func TestGetFileList(t *testing.T) {
	packListFn := func(root string) (*FileList, error) {
		var list FileList
		strip := filepath.Dir(filepath.Clean(root)) + "/"
		if strings.HasSuffix(root, "/") {
			strip = filepath.Clean(root) + "/"
		}
		fmt.Printf("filepath.Clean(root): %v\n", filepath.Clean(root))

		list.Strip = strip

		for f := range getFileList(root) {
			list.Files = append(list.Files, f)
			list.TotalSize += f.GetSize()
		}

		return &list, nil
	}

	t.Run("testGetLocal", func(t *testing.T) {
		fileList, err := packListFn("./")
		assert.Nil(t, err)
		t.Log(plog.Jsonify(fileList))
	})

	t.Run("testGetContent", func(t *testing.T) {
		fileList, err := packListFn("./content")
		assert.Nil(t, err)
		t.Log(plog.Jsonify(fileList))
	})

	t.Run("testGetAbsFilesync", func(t *testing.T) {
		fileList, err := packListFn("/Users/yong/programes/go/src/github.com/superwhys/remoteX/internal/filesync")
		assert.Nil(t, err)
		t.Log(plog.Jsonify(fileList))
	})

}
