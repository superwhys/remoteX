package sender

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/stretchr/testify/assert"
	"github.com/superwhys/remoteX/internal/filesync/pb"
)

func TestGetFileList(t *testing.T) {
	packListFn := func(root string) (*pb.FileList, error) {
		st := &SendTransfer{}
		var list pb.FileList
		strip := filepath.Dir(filepath.Clean(root)) + "/"
		if strings.HasSuffix(root, "/") {
			strip = filepath.Clean(root) + "/"
		}
		fmt.Printf("filepath.Clean(root): %v\n", filepath.Clean(root))
		
		list.Strip = strip
		
		for f := range st.GetFileList(root) {
			list.Files = append(list.Files, f)
			list.TotalSize += f.GetEntry().GetSize()
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
