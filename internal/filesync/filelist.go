package filesync

import (
	"io/fs"
	"iter"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-puzzles/puzzles/plog"
)

func (f *FileList) Sort() {
	sort.Slice(f.Files, func(i, j int) bool {
		return f.Files[i].Wpath < f.Files[j].Wpath
	})
}

func getFileList(root string) iter.Seq[*FileBase] {
	ch := make(chan *FileBase)
	go func() {
		strip := filepath.Dir(filepath.Clean(root)) + "/"
		if strings.HasSuffix(root, "/") {
			strip = filepath.Clean(root) + "/"
		}

		filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
			if err != nil {
				plog.Errorf("filepath walk error: %v", err)
				return err
			}

			if path == root {
				return nil
			}

			name := strings.TrimPrefix(path, strip)
			if name == root {
				name = "."
			}

			size := info.Size()

			f := &FileBase{
				Path:    name,
				Regular: info.Mode().IsRegular(),
				Wpath:   name,
				Size:    size,
			}
			ch <- f

			return nil
		})

		close(ch)
	}()

	return func(yield func(*FileBase) bool) {
		for f := range ch {
			if !yield(f) {
				break
			}
		}
	}
}
