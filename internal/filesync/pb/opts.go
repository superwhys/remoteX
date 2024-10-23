package pb

func (opts *SyncOpts) SetDefault() *SyncOpts {
	if opts == nil {
		opts = &SyncOpts{}
	}

	return opts
}
