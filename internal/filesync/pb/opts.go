package pb

func (m *SyncOpts) Merge(opt *SyncOpts) {
	if opt == nil {
		return
	}

	if m.DryRun == false && opt.DryRun != false {
		m.DryRun = opt.DryRun
	}

	if m.Whole == false && opt.Whole != false {
		m.Whole = opt.Whole
	}

	return
}

func (m *SyncOpts) SetDefault() *SyncOpts {
	if m == nil {
		m = &SyncOpts{}
	}

	return m
}
