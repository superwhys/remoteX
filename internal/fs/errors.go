package fs

type FsError uint64

const (
	ErrUnknown FsError = iota
	ErrDirPathNotFound
)

func (e FsError) Code() uint64 {
	return uint64(e)
}
func (e FsError) String() string {
	switch e {
	case ErrUnknown:
		return "unknown"
	case ErrDirPathNotFound:
		return "path not found"
	default:
		return "unknown"
	}
}
func (e FsError) Error() string {
	return e.String()
}
