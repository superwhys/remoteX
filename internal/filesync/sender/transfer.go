package sender

import (
	"context"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync"
	"github.com/superwhys/remoteX/internal/filesync/file"
	"github.com/superwhys/remoteX/internal/filesync/match"
	"github.com/superwhys/remoteX/internal/filesync/opts"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type sendTransfer struct {
	opts *opts.SyncOpt
	rw   protoutils.ProtoMessageReadWriter
}

func (st *sendTransfer) sendFileList(ctx context.Context, root string) (*filesync.FileList, error) {
	var fileList filesync.FileList

	strip := filepath.Dir(filepath.Clean(root)) + "/"
	if strings.HasSuffix(root, "/") {
		strip = filepath.Clean(root) + "/"
	}

	fileList.Strip = strip

	for f := range filesync.GetFileList(root) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		fileList.Files = append(fileList.Files, f)
		fileList.TotalSize += f.GetSize()

		if err := st.rw.WriteMessage(f); err != nil {
			return nil, errors.Wrapf(err, "write file: %s", f.GetPath())
		}
	}

	fileList.Sort()

	if err := st.rw.WriteMessage(&filesync.FileBase{IsEnd: true}); err != nil {
		return nil, errors.Wrap(err, "write end")
	}

	return &fileList, nil
}

func (st *sendTransfer) receiveFileIdx() (*filesync.FileIdx, error) {
	var fileIdx filesync.FileIdx
	if err := st.rw.ReadMessage(&fileIdx); err != nil {
		return nil, errors.Wrap(err, "read idx")
	}
	return &fileIdx, nil
}

func (st *sendTransfer) receiveHeadSum() (*filesync.HashHead, error) {
	var head filesync.HashHead
	if err := st.rw.ReadMessage(&head); err != nil {
		return nil, errors.Wrap(err, "read head")
	}

	if head.GetCheckSumCount() == 0 {
		// need whole file
		return &head, nil
	}

	head.Hashs = make([]*filesync.HashBuf, head.GetCheckSumCount())
	for i := int64(0); i < head.GetCheckSumCount(); i++ {
		var hashBuf filesync.HashBuf
		if err := st.rw.ReadMessage(&hashBuf); err != nil {
			return nil, errors.Wrapf(err, "read hashBuf: %v", i)
		}

		head.Hashs[i] = &hashBuf
	}
	return &head, nil
}

func (st *sendTransfer) sendFile(ctx context.Context, blockLength, fileSize int64, srcPath string) error {
	srcFile, err := file.OpenFile(srcPath)
	if err != nil {
		return errors.Wrapf(err, "open file: %s", srcPath)
	}
	defer srcFile.Close()

	var (
		offset int64
		l      = blockLength
	)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if offset+blockLength > fileSize {
			l = fileSize - offset
		}

		b, err := file.ReadFileAtOffset(srcFile, offset, l)
		if err != nil {
			return errors.Wrapf(err, "read file at offset: %d", offset)
		}

		if err := st.rw.WriteMessage(&filesync.FileChunk{Data: b}); err != nil {
			return errors.Wrapf(err, "write match chunk: %s", srcPath)
		}

		offset += blockLength
		if offset >= fileSize {
			break
		}
	}

	return st.rw.WriteMessage(&filesync.FileChunk{IsEnd: true})
}

func (st *sendTransfer) hashMatch(ctx context.Context, head *filesync.HashHead, srcPath string) error {
	srcFile, err := file.OpenFile(srcPath)
	if err != nil {
		return errors.Wrapf(err, "open file: %s", srcPath)
	}
	defer srcFile.Close()

	matchIter, err := match.HashMatch(ctx, head, srcFile)
	if err != nil {
		return errors.Wrapf(err, "hash match: %s", srcPath)
	}

	for matchChunk, err := range matchIter {
		if err != nil {
			return errors.Wrapf(err, "iter hash match: %s", srcPath)
		}

		if err := st.rw.WriteMessage(matchChunk); err != nil {
			return errors.Wrapf(err, "write match chunk: %s", srcPath)
		}
	}

	return st.rw.WriteMessage(&filesync.FileChunk{IsEnd: true})
}
