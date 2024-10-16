package filesync

type SyncOpt struct {
	// DryRun 模拟运行，不会进行实际的文件传输
	DryRun bool
}

func (opts *SyncOpt) UpdateOpt() {
}
