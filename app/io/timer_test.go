package io

import (
	"goboy/app/joypad"
	"goboy/app/ppu"
	"testing"
)

func TestTimerTick(t *testing.T) {
	timer := &Timer{}
	
	// DIV should increment every tick
	timer.Tick()
	if timer.div != 1 {
		t.Errorf("Expected div=1, got %d", timer.div)
	}

	// DIV Read should return top 8 bits
	timer.div = 0x1234
	if timer.ReadDIV() != 0x12 {
		t.Errorf("Expected ReadDIV=0x12, got %02X", timer.ReadDIV())
	}

	// Writing to DIV should reset it
	timer.WriteDIV(0x55)
	if timer.div != 0 {
		t.Errorf("Expected div=0 after write, got %d", timer.div)
	}
}

func TestTIMAIncrement(t *testing.T) {
	interruptRequested := false
	timer := &Timer{
		RequestInterrupt: func(bit uint8) {
			if bit == 2 {
				interruptRequested = true
			}
		},
	}

	// Enable timer with 4096Hz frequency (TAC bit 0, 1 = 00, bit 2 = 1)
	timer.WriteTAC(0x04)
	
	// 4096Hz means increment when bit 9 of div falls from 1 to 0
	// bit 9 is 0x0200
	
	timer.div = 0x01FF
	timer.Tick() // div becomes 0x0200, bit 9 becomes 1. lastAndResult becomes true.
	
	if timer.tima != 0 {
		t.Errorf("TIMA should not increment yet, got %d", timer.tima)
	}

	timer.div = 0x03FF
	timer.Tick() // div becomes 0x0400, bit 9 becomes 0. lastAndResult becomes false.
	// Falling edge detected, TIMA should increment.
	
	if timer.tima != 1 {
		t.Errorf("Expected TIMA=1, got %d", timer.tima)
	}

	if interruptRequested {
		t.Error("Interrupt should not be requested yet")
	}
}

func TestTIMAOverflow(t *testing.T) {
	interruptRequested := false
	timer := &Timer{
		tima: 0xFF,
		tma:  0x42,
		tac:  0x04, // Enabled
		RequestInterrupt: func(bit uint8) {
			if bit == 2 {
				interruptRequested = true
			}
		},
	}
	
	// Force an increment
	timer.div = 0x01FF
	timer.Tick() // div becomes 0x0200, bit 9 becomes 1. lastAndResult becomes true.
	
	timer.div = 0x03FF
	timer.Tick() // div becomes 0x0400, bit 9 becomes 0. lastAndResult becomes false.
	// TIMA becomes 0, overflowPending becomes true, overflowCycles = 4
	
	if timer.tima != 0 {
		t.Errorf("Expected TIMA=0, got %02X", timer.tima)
	}
	if interruptRequested {
		t.Error("Interrupt should not be requested immediately")
	}

	// 4 ticks later
	timer.Tick() // 3
	timer.Tick() // 2
	timer.Tick() // 1
	timer.Tick() // 0 - Reload and Interrupt
	
	if timer.tima != 0x42 {
		t.Errorf("Expected TIMA=0x42, got %02X", timer.tima)
	}
	if !interruptRequested {
		t.Error("Interrupt should be requested after 4 cycles")
	}
}

func TestIO(t *testing.T) {
	p := ppu.CreatePPU()
	p.OnStat = func() {}
	p.OnVBlank = func() {}
	
	j := joypad.CreateJoypad()
	timer := &Timer{}
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
	io.Write(DIV_ADDRESS, 0x55)
	if timer.div != 0 {
		t.Errorf("DIV write through IO failed, got %d", timer.div)
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
