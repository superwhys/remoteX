package filesync

import (
	"context"
	"os"
	"path/filepath"

	"github.com/go-puzzles/puzzles/plog"
	"github.com/go-puzzles/puzzles/putils"
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/internal/filesystem"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type fileSyncer struct {
	fs filesystem.FileSystem
}

func NewFileSyncer() *fileSyncer {
	return &fileSyncer{
		fs: filesystem.NewBasicFileSystem(),
	}
}

func (f *fileSyncer) Statistic(r *pb.SyncResp, file *pb.FileBase) *pb.SyncResp {
	if r == nil {
		r = &pb.SyncResp{}
	}

	transFile := &pb.SyncFile{
		Name: file.GetEntry().GetName(),
		Size: file.GetEntry().GetSize(),
		Type: file.GetEntry().GetType(),
	}
	r.Total++
	r.TotalSize += file.GetEntry().GetSize()
	r.ActualSendBytes = file.ActualSend
	r.Files = append(r.Files, transFile)

	return r
}

func (f fileSyncer) SendFiles(ctx context.Context, rw protoutils.ProtoMessageReadWriter, path string, opts *pb.SyncOpts) (resp *pb.SyncResp, err error) {
	ctx = plog.With(ctx, "Send")
	opts = opts.SetDefault()

	if err := rw.WriteMessage(opts); err != nil {
		return nil, errors.Wrap(err, "sendOpts")
	}

	resp = &pb.SyncResp{
		Files:      make([]*pb.SyncFile, 0, 10),
		ErrorFiles: make([]*pb.ErrorFile, 0, 10),
	}

	st := NewSendTransfer(rw)
	fileListIter := f.sendFileList(ctx, path, rw)
	for file, err := range fileListIter {
		if err != nil {
			plog.Errorc(ctx, "send fileList error: %v", err)
			return nil, errors.Wrap(err, "send fileList")
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if file.GetIsEnd() {
			break
		}

		cCtx := plog.With(ctx, "file", file.GetEntry().GetWpath())
		plog.Debugc(cCtx, "process file")

		if opts.DryRun {
			resp = f.Statistic(resp, file)
			continue
		}

		switch file.GetEntry().GetType() {
		case filesystem.EntryTypeDir:
		case filesystem.EntryTypeFile:
			err = st.Transfer(cCtx, file, opts)
			if err != nil {
				plog.Errorc(cCtx, "send file error: %v", err)
				resp.ErrorFiles = append(resp.ErrorFiles, &pb.ErrorFile{
					Name:    file.GetEntry().GetPath(),
					Message: err.Error(),
				})
				continue
			}
		default:
			continue
		}

		resp = f.Statistic(resp, file)

		ack := &pb.FileSyncAck{}
		err = rw.ReadMessage(ack)
		if err != nil {
			plog.Errorc(cCtx, "read file sync ack error: %v", err)
			return nil, errors.Wrap(err, "readFileSyncAck")
		}

		plog.Debugc(cCtx, "read sync file ack: %v", ack)
		plog.Infoc(cCtx, "transfer success.")
	}

	err = rw.ReadMessage(&pb.FileIdx{})
	if err != nil {
		plog.Errorc(ctx, "read done fileIdx error: %v", err)
		return nil, errors.Wrap(err, "readDoneFileIdx")
	}

	return resp, nil
}

func (f fileSyncer) ReceiveFiles(ctx context.Context, rw protoutils.ProtoMessageReadWriter, dest string, opts *pb.SyncOpts) (resp *pb.SyncResp, err error) {
	ctx = plog.With(ctx, "Receive")

	if err := f.checkDest(dest); err != nil {
		return nil, errors.Wrap(err, "checkDesk")
	}

	opts = opts.SetDefault()
	remoteOpts := new(pb.SyncOpts)
	if err := rw.ReadMessage(remoteOpts); err != nil {
		return nil, err
	}
	opts.Merge(remoteOpts)

	resp = &pb.SyncResp{
		Files:      make([]*pb.SyncFile, 0, 10),
		ErrorFiles: make([]*pb.ErrorFile, 0, 10),
	}

	rt := NewReceiverTransfer(rw, dest)
	fileListIter := f.receiveFileList(ctx, rw)
	for file, err := range fileListIter {
		if err != nil {
			plog.Errorc(ctx, "receive fileList error: %v", err)
			return nil, errors.Wrap(err, "receive fileList")
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if file.GetIsEnd() {
			break
		}

		cCtx := plog.With(ctx, "file", file.GetEntry().GetWpath())
		plog.Debugc(cCtx, "process file")

		if opts.DryRun {
			resp = f.Statistic(resp, file)
			continue
		}

		switch file.GetEntry().GetType() {
		case filesystem.EntryTypeDir:
			t := filepath.Join(dest, file.GetEntry().GetWpath())
			if !putils.FileExists(t) {
				err = os.MkdirAll(t, 0755)
				if err != nil {
					rw.WriteMessage(&pb.FileIdx{Idx: -1})
					return nil, errors.Wrapf(err, "create dir: %s", t)
				}
			}
		case filesystem.EntryTypeFile:
			err = rt.Transfer(cCtx, file, opts)
			if err != nil {
				plog.Errorc(cCtx, "receive file error: %v", err)
				rw.WriteMessage(&pb.FileIdx{Idx: -1})
				return nil, errors.Wrapf(err, "receive file: %s", file.GetEntry().GetPath())
			}
		default:
			continue
		}

		resp = f.Statistic(resp, file)

		err = rw.WriteMessage(&pb.FileSyncAck{
			Success:     true,
			ReceiveSize: rt.TransferSize(),
			Idx:         file.GetIdx(),
		})
		if err != nil {
			plog.Errorc(cCtx, "write file sync ack error: %v", err)
			return nil, errors.Wrap(err, "writeFileSyncAck")
		}
		plog.Infoc(cCtx, "transfer success.")
	}

	err = rw.WriteMessage(&pb.FileIdx{Idx: -1})
	if err != nil {
		plog.Errorc(ctx, "write done fileIdx error: %v", err)
		return nil, errors.Wrap(err, "writeDoneFileIdx")
	}

	return resp, nil
}
