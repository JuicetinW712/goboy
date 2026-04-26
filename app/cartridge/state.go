package cartridge

type CartridgeState struct {
	Ram      []uint8
	Title    string
	Filename string

	MbcType MBCType
	Mbc     any
}

func (c *Cartridge) GetState() CartridgeState {
	return CartridgeState{
		Ram:      c.ram,
		Title:    c.title,
		Filename: c.filename,
		MbcType:  c.mbcType,
		Mbc:      c.mbc.GetState(),
	}
}

func (c *Cartridge) LoadState(state CartridgeState) error {
	copy(c.ram, state.Ram)
	return c.mbc.LoadState(state.Mbc)
}
