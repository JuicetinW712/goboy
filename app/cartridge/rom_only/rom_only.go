package rom_only

type ROMOnlyMBC struct {
	rom []uint8
}

func CreateROMOnlyMBC(rom []uint8) *ROMOnlyMBC {
	return &ROMOnlyMBC{
		rom: rom,
	}
}

func (mbc *ROMOnlyMBC) Read(address uint16) uint8 {
	return mbc.rom[address]
}

// ROM Only cartridges are not writeable
func (mbc *ROMOnlyMBC) Write(address uint16, value uint8) {}
