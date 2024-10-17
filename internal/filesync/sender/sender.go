package sender

import (
	"context"
	"path/filepath"
	"strings"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/file"
	"github.com/superwhys/remoteX/internal/filesync/match"
	"github.com/superwhys/remoteX/internal/filesync/opts"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type SendTransfer struct {
	Opts *opts.SyncOpt
	Rw   protoutils.ProtoMessageReadWriter
}

func (st *SendTransfer) SendFileList(ctx context.Context, root string) (*pb.FileList, error) {
	var fileList pb.FileList
	
	strip := filepath.Dir(filepath.Clean(root)) + "/"
	if strings.HasSuffix(root, "/") {
		strip = filepath.Clean(root) + "/"
	}
	
	fileList.Strip = strip
	
	for f := range file.GetFileList(strip, root) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		fileList.Files = append(fileList.Files, f)
		fileList.TotalSize += f.GetSize()
		
		if err := st.Rw.WriteMessage(f); err != nil {
			return nil, errors.Wrapf(err, "write file: %s", f.GetName())
		}
	}
	
	fileList.Sort()
	
	if err := st.Rw.WriteMessage(&pb.FileBase{IsEnd: true}); err != nil {
		return nil, errors.Wrap(err, "write end")
	}
	
	return &fileList, nil
}

func (st *SendTransfer) ReceiveFileIdx() (*pb.FileIdx, error) {
	var fileIdx pb.FileIdx
	if err := st.Rw.ReadMessage(&fileIdx); err != nil {
		return nil, errors.Wrap(err, "read idx")
	}
	return &fileIdx, nil
}

func (st *SendTransfer) ReceiveHeadSum(ctx context.Context) (*pb.HashHead, error) {
	var head pb.HashHead
	if err := st.Rw.ReadMessage(&head); err != nil {
		return nil, errors.Wrap(err, "read head")
	}
	
	if head.GetCheckSumCount() == 0 {
		// need whole file
		return &head, nil
	}
	
	head.Hashs = make([]*pb.HashBuf, head.GetCheckSumCount())
	for i := int64(0); i < head.GetCheckSumCount(); i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		var hashBuf pb.HashBuf
		if err := st.Rw.ReadMessage(&hashBuf); err != nil {
			return nil, errors.Wrapf(err, "read hashBuf: %v", i)
		}
		
		head.Hashs[i] = &hashBuf
	}
	return &head, nil
}

func (st *SendTransfer) SendFile(ctx context.Context, blockLength, fileSize int64, srcPath string) error {
	plog.Debugf("start send whole file: %v", srcPath)
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
		
		if err := st.Rw.WriteMessage(&pb.FileChunk{Data: b}); err != nil {
			return errors.Wrapf(err, "write match chunk: %s", srcPath)
		}
		
		offset += blockLength
		if offset >= fileSize {
			break
		}
	}
	
	plog.Debugf("send whole file: %v done", srcPath)
	return st.Rw.WriteMessage(&pb.FileChunk{IsEnd: true})
}

func (st *SendTransfer) HashMatch(ctx context.Context, head *pb.HashHead, srcPath string) error {
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
		
		if err := st.Rw.WriteMessage(matchChunk); err != nil {
			return errors.Wrapf(err, "write match chunk: %s", srcPath)
		}
	}
	
	return st.Rw.WriteMessage(&pb.FileChunk{IsEnd: true})
}
