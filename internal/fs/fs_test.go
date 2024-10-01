package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListDir(t *testing.T) {
	basicFs := NewBasicFileSystem()

	entries, err := basicFs.List("./content")
	if !assert.Nil(t, err) {
		return
	}

	t.Log(entries)
}
