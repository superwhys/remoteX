package filesync

import (
	"context"
	"os"
	"path/filepath"
	
	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/opts"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/internal/filesync/receiver"
	"github.com/superwhys/remoteX/internal/filesync/sender"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

func SendFiles(ctx context.Context, rw protoutils.ProtoMessageReadWriter, path string, opts *opts.SyncOpt) (err error) {
	opts = opts.SetDefault()
	
	st := &sender.SendTransfer{
		Opts: opts,
		Rw:   rw,
	}
	
	fileList, err := st.SendFileList(ctx, path)
	if err != nil {
		return errors.Wrap(err, "sendFileList")
	}
	
	for {
		// receive client process file idx
		fileIdx, err := st.ReceiveFileIdx()
		if err != nil {
			return errors.Wrap(err, "receiveFileIdx")
		}
		if fileIdx.GetIdx() == -1 {
			break
		}
		
		file := fileList.GetFiles()[fileIdx.GetIdx()]
		plog.Debugf("receive file idx %d, file: %s", fileIdx.GetIdx(), file.GetEntry().GetName())
		
		if opts.DryRun {
			continue
		}
		
		// receive client file sums
		plog.Debugf("start receive file sums")
		head, err := st.ReceiveHeadSum(ctx)
		if err != nil {
			return errors.Wrap(err, "receiveHeadSum")
		}
		plog.Debugf("receive file sums count: %d", len(head.GetHashs()))
		
		// transfer file
		plog.Debugf("start transfer file: %v", file.GetEntry().GetName())
		srcPath := filepath.Join(fileList.GetStrip(), file.GetEntry().GetName())
		if len(head.GetHashs()) == 0 {
			plog.Debugf("head hashs is empty, need whole file")
			
			err = st.SendFile(ctx, head.GetBlockLength(), file.GetEntry().GetSize(), srcPath)
		} else {
			err = st.HashMatch(ctx, head, srcPath)
		}
		
		if err != nil {
			return errors.Wrapf(err, "transfer file: %s", srcPath)
		}
		
		plog.Debugf("transfer file success: %v", file.GetEntry().GetName())
	}
	
	return nil
}

func ReceiveFile(ctx context.Context, rw protoutils.ProtoMessageReadWriter, dest string, opts *opts.SyncOpt) error {
	opts = opts.SetDefault()
	
	rt := &receiver.ReceiveTransfer{
		Opts: opts,
		Dest: dest,
		Rw:   rw,
	}
	// receive fileList
	fileList, err := rt.ReceiveFileList(ctx)
	if err != nil {
		return errors.Wrap(err, "receiveFileList")
	}
	fileCnt := len(fileList.Files)
	plog.Debugf("receive file list, count: %d", fileCnt)
	
	info, err := os.Stat(dest)
	if err != nil && !os.IsNotExist(err) {
		return errors.Wrap(err, "stat")
	}
	
	if info == nil {
		// dest not exists
		// if multiple files are received, dest must be a folder
		// if just has one file and dest is not exists, it will be treated as a file
		if fileCnt > 1 {
			rt.DestIsDir = true
		}
	} else {
		rt.DestIsDir = info.IsDir()
	}
	
	if fileCnt > 1 && !rt.DestIsDir {
		return errors.New("dest is not a directory")
	}
	
	if rt.DestIsDir && info == nil {
		if err = os.MkdirAll(dest, 0755); err != nil {
			return errors.Wrapf(err, "mkdir： %s", dest)
		}
	}
	
	for idx, f := range fileList.Files {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		if err := rt.Rw.WriteMessage(&pb.FileIdx{Idx: int64(idx)}); err != nil {
			return errors.Wrapf(err, "sendFileIdx: %s", f.GetEntry().GetName())
		}
		plog.Debugf("write file idx: %d, file: %s", idx, f.GetEntry().GetName())
		
		if opts.DryRun {
			continue
		}
		
		plog.Debugf("start send file hash: %s", f.GetEntry().GetName())
		// generate each file sum and send
		fileExists, err := rt.CalcFileHashAndSend(ctx, f)
		if err != nil {
			return errors.Wrapf(err, "calculate file hash: %s", f.GetEntry().GetName())
		}
		plog.Debugf("send file hash success: %s", f.GetEntry().GetName())
		
		// receive server file match chunk
		plog.Debugf("start transfer file: %s", f.GetEntry().GetName())
		if err := rt.TransferFile(ctx, fileExists, dest, f); err != nil {
			return errors.Wrap(err, "transferFile")
		}
		plog.Debugf("transfer file success: %s", f.GetEntry().GetName())
	}
	
	return rt.Rw.WriteMessage(&pb.FileIdx{Idx: -1})
}
