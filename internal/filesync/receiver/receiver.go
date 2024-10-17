package receiver

import (
	"context"
	"iter"
	"os"
	"path/filepath"
	
	"github.com/pkg/errors"
	"github.com/superwhys/remoteX/internal/filesync/common"
	"github.com/superwhys/remoteX/internal/filesync/file"
	"github.com/superwhys/remoteX/internal/filesync/hash"
	"github.com/superwhys/remoteX/internal/filesync/match"
	"github.com/superwhys/remoteX/internal/filesync/opts"
	"github.com/superwhys/remoteX/internal/filesync/pb"
	"github.com/superwhys/remoteX/pkg/protoutils"
)

type ReceiveTransfer struct {
	Opts      *opts.SyncOpt
	Dest      string
	DestIsDir bool
	Rw        protoutils.ProtoMessageReadWriter
}

func (rt *ReceiveTransfer) ReceiveFileList(ctx context.Context) (*pb.FileList, error) {
	var fileList pb.FileList
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		
		f := &pb.FileBase{}
		if err := rt.Rw.ReadMessage(f); err != nil {
			return nil, err
		}
		
		if f.IsEnd {
			break
		}
		
		fileList.Files = append(fileList.Files, f)
		if f.Regular {
			fileList.TotalSize += f.Size
		}
	}
	
	fileList.Sort()
	
	return &fileList, nil
}

func (rt *ReceiveTransfer) CalcFileHashAndSend(ctx context.Context, f *pb.FileBase) (exists bool, err error) {
	var local string
	if rt.DestIsDir {
		local = filepath.Join(rt.Dest, f.GetName())
	} else {
		local = rt.Dest
	}
	st, err := os.Stat(local)
	if os.IsNotExist(err) {
		// need whole file
		return false, rt.Rw.WriteMessage(&pb.HashHead{
			BlockLength: int64(common.BlockSize),
		})
	} else if err != nil {
		return false, err
	}
	
	in, err := file.OpenFile(local)
	if err != nil {
		return false, errors.Wrapf(err, "openFile: %s", local)
	}
	fileLen := st.Size()
	
	head := hash.CalcHashHead(fileLen)
	if err := rt.Rw.WriteMessage(head); err != nil {
		return false, err
	}
	
	for hb := range hash.CalcFileSubHash(head, fileLen, in.File()) {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		default:
		}
		
		if err := rt.Rw.WriteMessage(hb); err != nil {
			return false, errors.Wrap(err, "sendHashBuf")
		}
	}
	
	return true, nil
}

func (rt *ReceiveTransfer) receiveFileChunkIter(ctx context.Context) (matchIter iter.Seq2[*pb.FileChunk, error]) {
	return func(yield func(*pb.FileChunk, error) bool) {
		for {
			var fileChunk pb.FileChunk
			if err := rt.Rw.ReadMessage(&fileChunk); err != nil {
				yield(nil, errors.Wrapf(err, "readFileChunk"))
				return
			}
			
			if fileChunk.GetIsEnd() || !yield(&fileChunk, nil) {
				return
			}
			
			select {
			case <-ctx.Done():
				yield(nil, ctx.Err())
				return
			default:
			}
		}
	}
}

func (rt *ReceiveTransfer) TransferFile(ctx context.Context, fileExists bool, dest string, fb *pb.FileBase) (err error) {
	var (
		target     *file.File
		targetPath string
	)
	if rt.DestIsDir {
		targetPath = filepath.Join(dest, fb.GetName())
	} else {
		targetPath = dest
	}
	
	if !fileExists {
		target, err = file.CreateFile(targetPath)
	} else {
		target, err = file.OpenFile(targetPath)
	}
	if err != nil {
		return errors.Wrap(err, "openFile")
	}
	
	matchIter := rt.receiveFileChunkIter(ctx)
	return match.SyncFile(ctx, matchIter, target)
}
