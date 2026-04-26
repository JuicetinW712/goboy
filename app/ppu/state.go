package ppu

import (
	"image"
)

type PPUState struct {
	VRAM              [VRAM_SIZE]uint8
	OAM               [OAM_SIZE]uint8
	Image             *image.RGBA
	Buffer            *image.RGBA
	LCDC              uint8
	STAT              uint8
	LY                uint8
	LYC               uint8
	LowStatValue      bool
	SCX               uint8
	SCY               uint8
	WY                uint8
	WX                uint8
	WindowLineCounter uint8
	BGP               uint8
	OBP0              uint8
	OBP1              uint8
	DMA               uint8
	Dot               uint64
}

func (p *PPU) GetState() PPUState {
	return PPUState{
		VRAM:              p.VRAM,
		OAM:               p.OAM,
		Image:             p.Image,
		Buffer:            p.Buffer,
		LCDC:              p.LCDC,
		STAT:              p.STAT,
		LY:                p.LY,
		LYC:               p.LYC,
		LowStatValue:      p.lowStatValue,
		SCX:               p.SCX,
		SCY:               p.SCY,
		WY:                p.WY,
		WX:                p.WX,
		WindowLineCounter: p.windowLineCounter,
		BGP:               p.BGP,

		OBP0: p.OBP0,
		OBP1: p.OBP1,
		DMA:  p.DMA,
		Dot:  p.dot,
	}
}

func (p *PPU) LoadState(state *PPUState) {
	p.VRAM = state.VRAM
	p.OAM = state.OAM
	p.LCDC = state.LCDC
	p.STAT = state.STAT
	p.LY = state.LY
	p.LYC = state.LYC
	p.lowStatValue = state.LowStatValue
	p.SCX = state.SCX
	p.SCY = state.SCY
	p.WY = state.WY
	p.WX = state.WX
	p.windowLineCounter = state.WindowLineCounter
	p.BGP = state.BGP
	p.OBP0 = state.OBP0
	p.OBP1 = state.OBP1
	p.DMA = state.DMA
	p.dot = state.Dot

	// Copy image and buffer if they are provided in state
	if state.Image != nil {
		copy(p.Image.Pix, state.Image.Pix)
	}
	if state.Buffer != nil {
		copy(p.Buffer.Pix, state.Buffer.Pix)
	}
}

func NewPPUFromState(state *PPUState) *PPU {
	return &PPU{
		VRAM:         state.VRAM,
		OAM:          state.OAM,
		Image:        state.Image,
		Buffer:       state.Buffer,
		LCDC:         state.LCDC,
		STAT:         state.STAT,
		LY:           state.LY,
		LYC:          state.LYC,
		lowStatValue: state.LowStatValue,
		SCX:          state.SCX,
		SCY:          state.SCY,
		WY:           state.WY,
		WX:           state.WX,
		BGP:          state.BGP,
		OBP0:         state.OBP0,
		OBP1:         state.OBP1,
		DMA:          state.DMA,
		dot:          state.Dot,
	}
}
