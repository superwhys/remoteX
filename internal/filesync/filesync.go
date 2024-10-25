package filesync

import (
	"context"
	"path/filepath"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/internal/filesync/receiver"
	"github.com/superwhys/remoteX/internal/filesync/sender"
	"github.com/superwhys/remoteX/internal/filesystem"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

func SendFiles(ctx context.Context, rw protoutils.ProtoMessageReadWriter, path string, opts *pb.SyncOpts) (resp *pb.SyncResp, err error) {
	defer plog.Infof("Send Files: %v done", path)
	opts = opts.SetDefault()

	st := &sender.SendTransfer{
		Opts: opts,
		Rw:   rw,
	}

	if err := st.SendOpts(opts); err != nil {
		return nil, errors.Wrap(err, "sendOpts")
	}

	fileList, err := st.SendFileList(ctx, path)
	if err != nil {
		return nil, errors.Wrap(err, "sendFileList")
	}
	fileCnt := len(fileList.Files)

	resp = &pb.SyncResp{
		Files:      make([]*pb.SyncFile, 0, fileCnt),
		ErrorFiles: make([]*pb.ErrorFile, 0, 10),
	}

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// receive client process file idx
		fileIdx, err := st.ReceiveFileIdx()
		if err != nil {
			return nil, errors.Wrap(err, "receiveFileIdx")
		}
		if fileIdx.GetIdx() == -1 || fileIdx.GetIdx() == 0 || int(fileIdx.GetIdx()) > fileCnt {
			break
		}

		file := fileList.GetFiles()[fileIdx.GetIdx()-1]
		srcPath := filepath.Join(fileList.GetStrip(), file.GetEntry().GetWpath())

		if opts.DryRun {
			resp = st.Statistic(resp, file)
			continue
		}

		err = func(ctx context.Context) error {
			ctx = plog.With(ctx, "file", srcPath)
			plog.Debugf("receive file idx %d", fileIdx.GetIdx())

			// receive client file sums
			plog.Debugc(ctx, "start receive head sum")
			head, err := st.ReceiveHeadSum(ctx)
			if err != nil {
				plog.Errorc(ctx, "receive head sum error: %v", err)
				return errors.Wrap(err, "receiveHeadSum")
			}
			plog.Debugc(ctx, "receive head sum: %v", head.GetCheckSumCount())

			// transfer file
			if len(head.GetHashs()) == 0 {
				err = st.SendFile(ctx, file.GetEntry().GetSize(), srcPath)
			} else {
				err = st.HashMatch(ctx, head, srcPath)
			}

			if err != nil {
				plog.Errorc(ctx, "transfer error: %v", err)
				return errors.Wrap(err, "transferFile")
			}

			// statistic
			resp = st.Statistic(resp, file)
			plog.Infoc(ctx, "transfer(%d/%d) success.", fileIdx.GetIdx(), len(fileList.Files))

			return nil
		}(ctx)
		if err != nil {
			plog.Errorf("send file %s error: %v", srcPath, err)
			resp.ErrorFiles = append(resp.ErrorFiles, &pb.ErrorFile{
				Name:    file.GetEntry().GetPath(),
				Message: err.Error(),
			})
			continue
		}
	}

	return resp, nil
}

func ReceiveFile(ctx context.Context, rw protoutils.ProtoMessageReadWriter, dest string, opts *pb.SyncOpts) error {
	opts = opts.SetDefault()

	rt := &receiver.ReceiveTransfer{
		Opts: opts,
		Dest: dest,
		Rw:   rw,
	}

	err := rt.MergeRemoteOpts()
	if err != nil {
		return errors.Wrap(err, "mergeRemoteOpts")
	}
	opts = rt.Opts

	// receive fileList
	fileList, err := rt.ReceiveFileList(ctx)
	if err != nil {
		return errors.Wrap(err, "receiveFileList")
	}
	fileCnt := len(fileList.Files)
	plog.Debugf("receive file list, count: %d", fileCnt)

	if err := rt.CheckDesk(fileCnt, dest); err != nil {
		return errors.Wrap(err, "checkDesk")
	}

	for idx, f := range fileList.Files {
		idx = idx + 1
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if f.Entry.Type == filesystem.EntryTypeDir {
			continue
		}

		if err := rt.Rw.WriteMessage(&pb.FileIdx{Idx: int64(idx)}); err != nil {
			return errors.Wrapf(err, "sendFileIdx: %s", f.GetEntry().GetName())
		}

		if opts.DryRun {
			continue
		}

		var targetPath string
		if rt.DestIsDir {
			targetPath = filepath.Join(dest, f.GetEntry().GetWpath())
		} else {
			targetPath = dest
		}

		err = func(ctx context.Context) error {
			ctx = plog.With(ctx, "file", targetPath)

			// generate each file sum and send
			plog.Debugc(ctx, "start calc file hash")
			err := rt.CalcFileHashAndSend(ctx, targetPath)
			if err != nil {
				return errors.Wrapf(err, "calculate file hash: %s", targetPath)
			}
			plog.Debugc(ctx, "calc file hash success")

			// receive server file match chunk
			plog.Debugc(ctx, "start transfer file")
			if err := rt.TransferFile(ctx, targetPath); err != nil {
				return errors.Wrap(err, "transferFile")
			}
			plog.Infoc(ctx, "transfer(%d/%d) success.", idx, len(fileList.Files))
			return nil
		}(ctx)
		if err != nil {
			plog.Errorf("receive file %s error: %v", targetPath, err)
			rt.Rw.WriteMessage(&pb.FileIdx{Idx: -1})
			return errors.Wrapf(err, "receive file: %s", targetPath)
		}
	}

	return rt.Rw.WriteMessage(&pb.FileIdx{Idx: -1})
}
