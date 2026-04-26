package io

import (
	"goboy/app/joypad"
	"goboy/app/ppu"
	"goboy/app/timer"
	"testing"
)

func TestIO(t *testing.T) {
	p := ppu.CreatePPU()
	p.OnStat = func() {}
	p.OnVBlank = func() {}

	j := joypad.CreateJoypad()
	timer := &timer.Timer{}
	io := NewIO(timer, p, j)

	// Test IF
	io.Write(IF_ADDRESS, 0x01)
	if io.Read(IF_ADDRESS) != 0xE1 {
		t.Errorf("IF read failed, got %02X", io.Read(IF_ADDRESS))
	}

	io.RequestInterrupt(2)
	if io.Read(IF_ADDRESS) != 0xE5 {
		t.Errorf("IF request failed, got %02X", io.Read(IF_ADDRESS))
	}

	// Test PPU access through IO
	io.Write(LCDC_ADDRESS, 0x91)
	if p.LCDC != 0x91 {
		t.Errorf("LCDC write through IO failed, got %02X", p.LCDC)
	}
	if io.Read(LCDC_ADDRESS) != 0x91 {
		t.Errorf("LCDC read through IO failed, got %02X", io.Read(LCDC_ADDRESS))
	}

	// Test Timer access through IO
	// TODO - probably fails now since its only reading upper byte?
	io.Write(DIV_ADDRESS, 0x55)
	if timer.ReadDIV() != 0 {
		t.Errorf("DIV write through IO failed, got %d", timer.ReadDIV())
	}

	io.Write(TIMA_ADDRESS, 0xAA)
	if io.Read(TIMA_ADDRESS) != 0xAA {
		t.Errorf("TIMA access failed")
	}

	io.Write(TMA_ADDRESS, 0xBB)
	if io.Read(TMA_ADDRESS) != 0xBB {
		t.Errorf("TMA access failed")
	}

	io.Write(TAC_ADDRESS, 0x07)
	if io.Read(TAC_ADDRESS) != 0x07 {
		t.Errorf("TAC access failed")
	}

	// Test Joypad access through IO
	io.Write(JOYPAD_ADDRESS, 0x10) // Select buttons (bit 4 is 0)
	// sel = 0x10. res = j.buttons = 0x0F. return 0xC0 | 0x10 | 0x0F = 0xDF.
	if io.Read(JOYPAD_ADDRESS) != 0xDF {
		t.Errorf("Joypad read failed, got %02X", io.Read(JOYPAD_ADDRESS))
	}
}
