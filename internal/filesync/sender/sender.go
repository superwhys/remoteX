package sender

import (
	"context"
	"path/filepath"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/opts"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type FileSender interface {
	SendFiles(ctx context.Context, rw protoutils.ProtoMessageReadWriter, path string, opts *opts.SyncOpt) (err error)
}

type Sender struct{}

func (s *Sender) SendFiles(ctx context.Context, rw protoutils.ProtoMessageReadWriter, path string, opts *opts.SyncOpt) (err error) {
	st := &sendTransfer{
		opts: opts,
		rw:   rw,
	}

	fileList, err := st.sendFileList(ctx, path)
	if err != nil {
		return errors.Wrap(err, "sendFileList")
	}

	for {
		// receive client process file idx
		fileIdx, err := st.receiveFileIdx()
		if err != nil {
			return errors.Wrap(err, "receiveFileIdx")
		}
		if fileIdx.GetIdx() == -1 {
			break
		}

		file := fileList.GetFiles()[fileIdx.GetIdx()]
		plog.Infof("receive file idx %d, fileName: %s", fileIdx.GetIdx(), file.Name())
		if opts.DryRun {
			continue
		}

		// receive client file sums
		head, err := st.receiveHeadSum()
		if err != nil {
			return errors.Wrap(err, "receiveHeadSum")
		}
		plog.Debugf("receive head hash sun count: %d", len(head.GetHashs()))

		// transfer file
		srcPath := filepath.Join(fileList.GetStrip(), file.GetPath())
		if len(head.GetHashs()) == 0 {
			plog.Debugf("head hashs is empty, need whole file")
			// send whole file list
			err = st.sendFile(ctx, head.GetBlockLength(), file.GetSize(), srcPath)
		} else {
			err = st.hashMatch(ctx, head, srcPath)
		}

		if err != nil {
			return errors.Wrapf(err, "transfer file: %s", srcPath)
		}
	}

	return nil
}
