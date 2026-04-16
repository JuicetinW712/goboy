package rom_only

type ROMOnlyState struct{}

func (mbc *ROMOnlyMBC) GetState() any {
	return ROMOnlyState{}
}

func (mbc *ROMOnlyMBC) LoadState(state any) error { return nil }
