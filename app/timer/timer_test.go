package timer

import (
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
