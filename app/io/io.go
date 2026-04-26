package io

import (
	"goboy/app/joypad"
	"goboy/app/ppu"
	"goboy/app/timer"
)

const (
	IO_START uint16 = 0xFF00
	IO_END   uint16 = 0xFF7F
	IO_SIZE  uint16 = IO_END - IO_START + 1

	JOYPAD_ADDRESS uint16 = 0xFF00
	IF_ADDRESS     uint16 = 0xFF0F

	DIV_ADDRESS  uint16 = 0xFF04
	TIMA_ADDRESS uint16 = 0xFF05
	TMA_ADDRESS  uint16 = 0xFF06
	TAC_ADDRESS  uint16 = 0xFF07

	// PPU Registers
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
)

type IO struct {
	timer  *timer.Timer
	ppu    *ppu.PPU
	joypad *joypad.Joypad

	IF uint8
}

func NewIO(timer *timer.Timer, ppu *ppu.PPU, joypad *joypad.Joypad) *IO {
	return &IO{
		timer:  timer,
		ppu:    ppu,
		joypad: joypad,
	}
}

func (i *IO) Read(address uint16) uint8 {
	switch {
	case address == JOYPAD_ADDRESS:
		return i.joypad.Read()
	case address == IF_ADDRESS:
		return i.IF | 0xE0 // Top 3 bits always read as 1
	case address == DIV_ADDRESS:
		return i.timer.ReadDIV()
	case address == TIMA_ADDRESS:
		return i.timer.ReadTIMA()
	case address == TMA_ADDRESS:
		return i.timer.ReadTMA()
	case address == TAC_ADDRESS:
		return i.timer.ReadTAC()
	case address >= 0xFF40 && address <= 0xFF4B:
		return i.ppu.Read(address)
	}
	return 0xFF
}

func (i *IO) Write(address uint16, value uint8) {
	switch {
	case address == JOYPAD_ADDRESS:
		i.joypad.Write(value)
	case address == IF_ADDRESS:
		i.IF = value
	case address == DIV_ADDRESS:
		i.timer.WriteDIV(value)
	case address == TIMA_ADDRESS:
		i.timer.WriteTIMA(value)
	case address == TMA_ADDRESS:
		i.timer.WriteTMA(value)
	case address == TAC_ADDRESS:
		i.timer.WriteTAC(value)
	case address >= 0xFF40 && address <= 0xFF4B:
		i.ppu.Write(address, value)
	}
}

func (i *IO) RequestInterrupt(bit uint8) {
	i.IF |= (1 << bit)
}
