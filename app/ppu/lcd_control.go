package ppu

// LCDC
// 7 - LCD & PPU enable: 0 = Off; 1 = On
// 6 - Window tile map area: 0 = 9800–9BFF; 1 = 9C00–9FFF
// 5 - Window enable: 0 = Off; 1 = On
// 4 - BG & Window tile data area: 0 = 8800–97FF; 1 = 8000–8FFF
// 3 - BG tile map area: 0 = 9800–9BFF; 1 = 9C00–9FFF
// 2 - OBJ size: 0 = 8×8; 1 = 8×16
// 1 - OBJ enable: 0 = Off; 1 = On
// 0 - BG & Window enable / priority [Different meaning in CGB Mode]: 0 = Off; 1 = On

func (p *PPU) lcdEnabled() bool {
	return (p.LCDC >> 7) == 1
}

func (p *PPU) getWindowTileMap() uint16 {
	if (p.LCDC>>6)&1 == 1 {
		return 1 // 0x9C00
	} else {
		return 0 // 0x9800
	}
}

func (p *PPU) windowEnabled() bool {
	return (p.LCDC>>5)&1 == 1
}

// Controls where the tile data is gotten from (excluding objects/sprites)
// Controls how addressing works
// 1 - unsigned, tile idx 255 is 255 (0  255) and centered around 8000
// 0 - signed, tile index 255 is -1 (-128 to 127) and centered around 9000
func (p *PPU) unsignedTileAddressing() bool {
	return (p.LCDC>>4)&1 == 1
}

// Controls which tile map is used to arrange background for window
func (p *PPU) getBackgroundTileMap() uint16 {
	if (p.LCDC>>3)&1 == 1 {
		return 1 // 0x9C00
	} else {
		return 0 // 0x9800
	}
}

// Controls the size of all objects (1 tile or 2 stacked vertically)
func (p *PPU) getObjSize() uint16 {
	if (p.LCDC>>2)&1 == 1 {
		return 16 // 8 x 16
	} else {
		return 8 // 8 x 8
	}
}

// Set if sprites are enabled
func (p *PPU) objEnabled() bool {
	return (p.LCDC>>1)&1 == 1
}

// 1 -> 0 means background becomes white
func (p *PPU) backgroundAndWindowEnabled() bool {
	return (p.LCDC & 1) == 1
}
