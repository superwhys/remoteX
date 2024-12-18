// File:		filelist.go
// Created by:	Hoven
// Created on:	2024-12-16
//
// This file is part of the Example Project.
//
// (c) 2024 Example Corp. All rights reserved.

package filesync

import (
	"context"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/putils"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

func (f *fileSyncer) getFileList(strip, root string) iter.Seq[*pb.FileBase] {
	filter := func(path string, info fs.FileInfo) bool {
		return info.IsDir() || info.Mode().IsRegular()
	}

	return func(yield func(*pb.FileBase) bool) {
		walkIter, err := f.fs.WalkIter(root, filter)
		if err != nil {
			plog.Errorf("WalkIter error: %v", err)
			return
		}
		idx := 1
		for entry := range walkIter {
			f := &pb.FileBase{
				Idx:   int64(idx),
				Strip: strip,
				Entry: entry,
			}

			if !yield(f) {
				break
			}

			idx++
		}
	}
}

func (f *fileSyncer) sendFileList(ctx context.Context, root string, rw protoutils.ProtoMessageReadWriter) iter.Seq2[*pb.FileBase, error] {
	strip := filepath.Dir(filepath.Clean(root)) + "/"
	if strings.HasSuffix(root, "/") {
		strip = filepath.Clean(root) + "/"
	}

	return func(yield func(*pb.FileBase, error) bool) {
		for f := range f.getFileList(strip, root) {
			select {
			case <-ctx.Done():
				return
			default:
			}

			err := rw.WriteMessage(f)
			if err != nil && !yield(nil, errors.Wrapf(err, "write file: %s", f.GetEntry().GetName())) {
				return
			}

			if !yield(f, nil) {
				return
			}
		}

		df := &pb.FileBase{IsEnd: true}
		err := rw.WriteMessage(df)
		if err != nil {
			yield(nil, errors.Wrap(err, "write end"))
		}

		yield(df, nil)
	}
}

func (f *fileSyncer) receiveFileList(ctx context.Context, rw protoutils.ProtoMessageReadWriter) iter.Seq2[*pb.FileBase, error] {
	return func(yield func(*pb.FileBase, error) bool) {
		for {
			select {
			case <-ctx.Done():
				break
			default:
			}

			file := &pb.FileBase{}
			err := rw.ReadMessage(file)
			if err != nil && !yield(nil, errors.Wrap(err, "readFileBase")) {
				break
			}

			if file.IsEnd {
				break
			}

			if !yield(file, nil) {
				break
			}
		}
	}
}

func (f *fileSyncer) checkDest(dest string) error {
	if !putils.FileExists(dest) {
		return os.MkdirAll(dest, 0755)
	}

	return nil
}
