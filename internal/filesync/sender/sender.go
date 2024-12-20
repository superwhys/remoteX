package sender

import (
	"context"
	"io"
	"path/filepath"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/match"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/internal/filesystem"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type SendTransfer struct {
	Fs         filesystem.FileSystem
	Rw         protoutils.ProtoMessageReadWriter
	ActualSend int
}

func (st *SendTransfer) SendOpts(opts *pb.SyncOpts) error {
	return st.Rw.WriteMessage(opts)
}

func (st *SendTransfer) TransferSize() int64 {
	return int64(st.ActualSend)
}

func (st *SendTransfer) Transfer(ctx context.Context, file *pb.FileBase, opts *pb.SyncOpts) (err error) {
	srcPath := filepath.Join(file.GetStrip(), file.GetEntry().GetWpath())
	ctx = plog.With(ctx, "file", srcPath)

	// receive client file sums
	plog.Debugc(ctx, "start receive head sum")
	head, err := st.receiveHeadSum(ctx)
	if err != nil {
		plog.Errorc(ctx, "receive head sum error: %v", err)
		return errors.Wrap(err, "receiveHeadSum")
	}
	plog.Debugc(ctx, "receive head sum: %v", head.GetCheckSumCount())

	// transfer file
	if len(head.GetHashs()) == 0 {
		err = st.sendFile(ctx, file.GetEntry().GetSize(), srcPath)
	} else {
		err = st.hashMatch(ctx, head, srcPath)
	}

	if err != nil {
		plog.Errorc(ctx, "transfer error: %v", err)
		return errors.Wrap(err, "transferFile")
	}

	file.ActualSend = int64(st.ActualSend)
	return nil
}

func (st *SendTransfer) receiveHeadSum(ctx context.Context) (*pb.HashHead, error) {
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

func (st *SendTransfer) sendFile(ctx context.Context, fileSize int64, srcPath string) (err error) {
	plog.Debugf("start send whole file: %v", srcPath)
	defer func() {
		if err != nil {
			plog.Errorf("send whole file has error: %v", err)
			st.Rw.WriteMessage(&pb.FileChunk{IsEnd: true})
		}
	}()

	srcFile, err := st.Fs.OpenFile(srcPath)
	if err != nil {
		return errors.Wrapf(err, "open file: %s", srcPath)
	}
	defer srcFile.Close()

	var (
		blockLength int64 = 256 * 1024
	)

	buf := make([]byte, blockLength)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := srcFile.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		chunk := buf[:n]

		if err := st.Rw.WriteMessage(&pb.FileChunk{Data: chunk}); err != nil {
			return errors.Wrapf(err, "write match chunk: %s", srcPath)
		}

		st.ActualSend += len(chunk)
		plog.Debugf("send %d bytes, total %d bytes", len(chunk), st.ActualSend)
	}

	return st.Rw.WriteMessage(&pb.FileChunk{IsEnd: true})
}

func (st *SendTransfer) hashMatch(ctx context.Context, head *pb.HashHead, srcPath string) (err error) {
	plog.Debugf("start send chunk file: %v", srcPath)
	defer func() {
		if err != nil {
			plog.Errorf("send hash match file has error: %v", err)
			st.Rw.WriteMessage(&pb.FileChunk{IsEnd: true})
		}
	}()

	srcFile, err := st.Fs.OpenFile(srcPath)
	if err != nil {
		return errors.Wrapf(err, "open file: %s", srcPath)
	}
	defer srcFile.Close()

	fi, err := srcFile.Stat()
	if err != nil {
		return errors.Wrapf(err, "stat file: %s", srcPath)
	}

	if fi.Size() == 0 {
		return st.Rw.WriteMessage(&pb.FileChunk{IsEnd: true})
	}

	matchIter, err := match.HashMatch(ctx, head, srcFile, fi)
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
		st.ActualSend += len(matchChunk.GetData())
	}

	return st.Rw.WriteMessage(&pb.FileChunk{IsEnd: true})
}
