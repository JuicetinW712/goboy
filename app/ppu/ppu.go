package ppu

import (
	"image"
	"image/color"
)

// TODO - when moving from 1 to 0 on background+windowenabled, should make entire line 0
// TODO - Need to fix window disabling and renabling to split the window
// TODO - enforce ten sprites per scanline
// TODO - DMA should incur a cycle penalty

const (
	VRAM_SIZE uint16 = 0x2000
	OAM_SIZE  uint16 = 0xA0

	VRAM_START uint16 = 0x8000
	VRAM_END   uint16 = 0x9FFF

	OAM_START uint16 = 0xFE00
	OAM_END   uint16 = 0xFE9F

	LCDC_ADDRESS uint16 = 0xFF40
	STAT_ADDRESS uint16 = 0xFF41
	SCY_ADDRESS  uint16 = 0xFF42
	SCX_ADDRESS  uint16 = 0xFF43
	LY_ADDRESS   uint16 = 0xFF44
	LYC_ADDRESS  uint16 = 0xFF45
	DMA_ADDRESS  uint16 = 0xFF46
	BGP_ADDRESS  uint16 = 0xFF47
	OBP0_ADDRESS uint16 = 0xFF48
	OBP1_ADDRESS uint16 = 0xFF49
	WY_ADDRESS   uint16 = 0xFF4A
	WX_ADDRESS   uint16 = 0xFF4B

	TILE_PIXEL_WIDTH uint16 = 8
	TILES_PER_ROW    uint16 = 32
	TILE_BYTE_SIZE          = 16

	DOTS_PER_LINE = 456
)

var colors = [4]color.RGBA{
	{0xFF, 0xFF, 0xFF, 0xFF}, // White
	{0xCC, 0xCC, 0xCC, 0xFF}, // Light Gray
	{0x77, 0x77, 0x77, 0xFF}, // Dark Gray
	{0x00, 0x00, 0x00, 0xFF}, // Black
}

type PPU_Mode uint8

const (
	HBlank PPU_Mode = iota // End of line
	VBlank
	OAMScan
	Drawing
)

// Gameboy Sprite - 1 row is 2 bytes, 8 rows = 16 bytes
// Ex: 00 00 01 01 10 10 11 11 (Example row)
// 00110011   lower bit in row 1
// 00001111   upper bit in row 2

// https://elrindel.github.io/specifications-gameboy

type PPU struct {
	VRAM   [VRAM_SIZE]uint8
	OAM    [OAM_SIZE]uint8
	Image  *image.RGBA
	Buffer *image.RGBA

	// LCDC
	// 7 - LCD & PPU enable: 0 = Off; 1 = On
	// 6 - Window tile map area: 0 = 9800–9BFF; 1 = 9C00–9FFF
	// 5 - Window enable: 0 = Off; 1 = On
	// 4 - BG & Window tile data area: 0 = 8800–97FF; 1 = 8000–8FFF
	// 3 - BG tile map area: 0 = 9800–9BFF; 1 = 9C00–9FFF
	// 2 - OBJ size: 0 = 8×8; 1 = 8×16
	// 1 - OBJ enable: 0 = Off; 1 = On
	// 0 - BG & Window enable / priority [Different meaning in CGB Mode]: 0 = Off; 1 = On
	LCDC uint8

	// LCD Status Registers
	STAT uint8 // LCD status
	LY   uint8 // LCD Y, Current horizontal line (0-153 normally, 144-153 VBlank)
	LYC  uint8 // LY compare, when LYC = LY, LYC=LY flag set in STAT register (for interrupt)

	// Account for STAT blocking
	// only a 0 -> 1 transition should trigger STAT
	lowStatValue bool

	// Scrolling Registers

	// Background offset
	// 256 * 256 tile map but screen is 160 * 144
	// Defines offset for scrolling (wraps around)
	SCX uint8 // (x + SCX) % 256
	SCY uint8 // (y + SCY) % 256

	// Window offset
	// Tell which x and y coord on the screen the window starts
	WY                uint8 // y
	WX                uint8 // x + 7
	windowLineCounter uint8

	// Palettes
	BGP  uint8 // background palette
	OBP0 uint8 // object palette 1
	OBP1 uint8 // object palette 2

	// Quick sprite cancer
	DMA uint8

	// Keeping track of PPU mode
	// https://gbdev.io/pandocs/Rendering.html
	dot uint64

	OnVBlank func()
	OnStat   func()
	OnDMA    func(value uint8)
}

func CreatePPU() *PPU {
	rectangle := image.Rect(0, 0, 160, 144)
	rectangle2 := image.Rect(0, 0, 160, 144)

	return &PPU{
		Image:  image.NewRGBA(rectangle),
		Buffer: image.NewRGBA(rectangle2),
	}
}

// https://gbdev.io/pandocs/Rendering.html
func (p *PPU) Step(dots uint) {
	if !p.lcdEnabled() {
		p.dot = 0
		p.LY = 0
		p.setMode(HBlank)
		return
	}

	p.dot += uint64(dots)

	if p.dot >= 456 {
		p.LY += 1
		p.dot -= 456
		p.updateLYC()
		p.checkStatInterrupt()
	}

	// Reset line number once we reach end of screen
	if p.LY > 153 {
		p.LY = 0
		p.windowLineCounter = 0
	}

	// After we reach end of screen (144), we mode to vblank mode
	if p.LY < 144 {
		// Mode 2 -> 3 -> 0
		if p.dot < 80 {
			p.handleOAMScanMode()
		} else if p.dot < 80+172 {
			p.handleDrawingMode()
		} else {
			p.handleHBlankMode()
		}
	} else {
		p.handleVBlankMode()
	}
}

// Mode 2
func (p *PPU) handleOAMScanMode() {
	// At the beginning of OAMScan, need to fire interrupt
	// if set in STAT
	if p.getMode() != OAMScan {
		p.setMode(OAMScan)
		p.checkStatInterrupt()
	}
}

// Mode 3
func (p *PPU) handleDrawingMode() {
	if p.getMode() != Drawing {
		p.setMode(Drawing)
		p.checkStatInterrupt()
	}
}

// Mode 0
func (p *PPU) handleHBlankMode() {
	if p.getMode() != HBlank {
		p.setMode(HBlank)
		p.drawLine()
		p.checkStatInterrupt()
	}
}

// Mode 1
func (p *PPU) handleVBlankMode() {
	if p.getMode() != VBlank {
		p.setMode(VBlank)

		// Copy from buffer to image once full frame rendered
		copy(p.Image.Pix, p.Buffer.Pix)

		p.OnVBlank()
		p.checkStatInterrupt()
	}
}

func (p *PPU) Read(addr uint16) uint8 {
	switch {
	case addr >= VRAM_START && addr <= VRAM_END:
		// VRAM blocked during Drawing (Mode 3)
		if p.getMode() == Drawing {
			return 0xFF
		}
		return p.VRAM[addr-VRAM_START]
	case addr >= 0xFE00 && addr <= 0xFE9F:
		// OAM blocked during OAMScan (Mode 2) and Drawing (Mode 3)
		mode := p.getMode()
		if mode == OAMScan || mode == Drawing {
			return 0xFF
		}
		return p.OAM[addr-OAM_START]
	case addr == LCDC_ADDRESS:
		return p.LCDC
	case addr == STAT_ADDRESS:
		return p.STAT
	case addr == SCY_ADDRESS:
		return p.SCY
	case addr == SCX_ADDRESS:
		return p.SCX
	case addr == LY_ADDRESS:
		return p.LY
	case addr == LYC_ADDRESS:
		return p.LYC
	case addr == BGP_ADDRESS:
		return p.BGP
	case addr == OBP0_ADDRESS:
		return p.OBP0
	case addr == OBP1_ADDRESS:
		return p.OBP1
	case addr == WY_ADDRESS:
		return p.WY
	case addr == WX_ADDRESS:
		return p.WX
	}
	return 0xFF
}

func (p *PPU) Write(addr uint16, val uint8) {
	switch {
	case addr >= VRAM_START && addr <= VRAM_END:
		// VRAM blocked during Drawing (Mode 3)
		if p.getMode() == Drawing {
			return
		}
		p.VRAM[addr-VRAM_START] = val
	case addr >= OAM_START && addr <= OAM_END:
		// OAM blocked during OAMScan (Mode 2) and Drawing (Mode 3)
		mode := p.getMode()
		if mode == OAMScan || mode == Drawing {
			return
		}
		p.OAM[addr-OAM_START] = val
	case addr == LCDC_ADDRESS:
		p.LCDC = val
	case addr == STAT_ADDRESS:
		p.STAT = (p.STAT & 0x07) | (val & 0xF8)
	case addr == SCY_ADDRESS:
		p.SCY = val
	case addr == SCX_ADDRESS:
		p.SCX = val
	case addr == LY_ADDRESS:
		// LY is read-only
	case addr == LYC_ADDRESS:
		p.LYC = val
	case addr == DMA_ADDRESS:
		p.DMA = val
		p.OnDMA(val)
	case addr == BGP_ADDRESS:
		p.BGP = val
	case addr == OBP0_ADDRESS:
		p.OBP0 = val
	case addr == OBP1_ADDRESS:
		p.OBP1 = val
	case addr == WY_ADDRESS:
		p.WY = val
	case addr == WX_ADDRESS:
		p.WX = val
	}
}

func (p *PPU) drawLine() {
	if p.backgroundAndWindowEnabled() {
		p.drawBackgroundLine()
	}

	if p.backgroundAndWindowEnabled() && p.windowEnabled() {
		p.drawWindowLine()
	}

	if p.objEnabled() {
		p.drawSpriteLine()
	}
}

func (p *PPU) drawBackgroundLine() {
	if !p.backgroundAndWindowEnabled() {
		return
	}

	// Get tile map to read from
	tileMapStart := uint16(0x9800)

	if p.getBackgroundTileMap() == 1 {
		tileMapStart = uint16(0x9C00)
	}

	// Get tile data
	y := uint16(p.LY + p.SCY) // Wraps around
	tilePixelY := y % TILE_PIXEL_WIDTH

	for i := range uint16(160) {
		x := (i + uint16(p.SCX)) % 256 // Wraps around
		tilePixelX := x % TILE_PIXEL_WIDTH

		// Get the tile number from the tile map
		tileX := x / TILE_PIXEL_WIDTH
		tileY := y / TILE_PIXEL_WIDTH

		tileMapAddr := tileMapStart + (tileY * TILES_PER_ROW) + tileX
		tileNum := p.VRAM[tileMapAddr-VRAM_START]

		// Get the pixel info from the tile
		tileDataStart := p.getTileAddress(tileNum) - VRAM_START
		lineDataLow := p.VRAM[tileDataStart+2*tilePixelY]
		lineDataHigh := p.VRAM[tileDataStart+2*tilePixelY+1]

		// Get 2 bit color info from the specified pixel in the line
		pixelDataLow := (lineDataLow >> (7 - tilePixelX)) & 1
		pixelDataHigh := (lineDataHigh >> (7 - tilePixelX)) & 1
		colorIdx := (pixelDataHigh << 1) | pixelDataLow

		// Set color data
		color := p.getColor(colorIdx, p.BGP)
		p.Buffer.SetRGBA(int(i), int(p.LY), color)
	}
}

func (p *PPU) drawWindowLine() {
	// Haven't reached window yet
	if p.LY < p.WY || p.WX >= 166 {
		return
	}

	// Get tile map to read from
	tileMapStart := uint16(0x9800)

	if p.getWindowTileMap() == 1 {
		tileMapStart = uint16(0x9C00)
	}

	// Get tile data
	windowY := uint16(p.windowLineCounter)
	tilePixelY := windowY % TILE_PIXEL_WIDTH

	startX := max(int(p.WX)-7, 0)

	for x := uint16(startX); x < 160; x++ {
		windowX := x - uint16(p.WX) + 7
		tilePixelX := windowX % TILE_PIXEL_WIDTH

		// Get the tile number from the tile map
		tileX := (windowX / TILE_PIXEL_WIDTH) % 32
		tileY := (windowY / TILE_PIXEL_WIDTH) % 32

		tileMapAddr := tileMapStart + (tileY * TILES_PER_ROW) + tileX
		tileNum := p.VRAM[tileMapAddr-VRAM_START]

		// Get the pixel info from the tile
		tileDataStart := p.getTileAddress(tileNum) - VRAM_START
		lineDataLow := p.VRAM[tileDataStart+2*tilePixelY]
		lineDataHigh := p.VRAM[tileDataStart+2*tilePixelY+1]

		// Get 2 bit color info from the specified pixel in the line
		pixelDataLow := (lineDataLow >> (7 - tilePixelX)) & 1
		pixelDataHigh := (lineDataHigh >> (7 - tilePixelX)) & 1
		colorIdx := (pixelDataHigh << 1) | pixelDataLow

		// Set color data
		color := p.getColor(colorIdx, p.BGP)
		p.Buffer.SetRGBA(int(x), int(p.LY), color)
	}

	p.windowLineCounter += 1
}

func (p *PPU) drawSpriteLine() {
	height := int16(p.getObjSize())

	for i := range 40 {
		offset := uint16(i * 4)
		y := int16(p.OAM[offset]) - 16
		x := int16(p.OAM[offset+1]) - 8
		tileID := p.OAM[offset+2]
		attr := p.OAM[offset+3]

		if int16(p.LY) < y || int16(p.LY) >= y+height {
			continue
		}

		palette := p.OBP0
		if (attr>>4)&1 == 1 {
			palette = p.OBP1
		}

		flipY := (attr>>6)&1 == 1
		flipX := (attr>>5)&1 == 1

		line := int16(p.LY) - y
		if flipY {
			line = height - 1 - line
		}

		tileAddr := uint16(tileID)*16 + uint16(line)*2
		byte1 := p.VRAM[tileAddr]
		byte2 := p.VRAM[tileAddr+1]

		for pixel := range int16(8) {
			pixelX := x + pixel
			if pixelX >= 160 {
				continue
			}

			bit := 7 - pixel
			if flipX {
				bit = pixel
			}

			colorIdx := (((byte2 >> bit) & 1) << 1) | ((byte1 >> bit) & 1)
			if colorIdx == 0 {
				continue // Transparent
			}

			// TODO: Check priority
			p.Buffer.SetRGBA(int(pixelX), int(p.LY), p.getColor(colorIdx, palette))
		}
	}
}

func (p *PPU) getTileAddress(tileNum uint8) uint16 {
	var tileDataAddr uint16

	if p.unsignedTileAddressing() {
		tileDataStart := 0x8000
		tileOffset := int(tileNum) * TILE_BYTE_SIZE
		tileDataAddr = uint16(tileDataStart + tileOffset)
	} else {
		tileDataStart := 0x9000
		tileOffset := int(int8(tileNum)) * TILE_BYTE_SIZE
		tileDataAddr = uint16(tileDataStart + tileOffset)
	}

	return tileDataAddr
}

// BGP, OGP0, OGP1 are palettes that you index into to get the color
// It's not fixed to allow for flexibility (e.g. want to make whole screen inverted colors
// only requires updating the BGP index)
func (p *PPU) getColor(colorIdx uint8, palette uint8) color.RGBA {
	idx := (palette >> (colorIdx * 2)) & 0x03
	return colors[idx]
}
