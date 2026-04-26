package mbc1

import (
	"testing"
)

func TestMBC1(t *testing.T) {
	// Create a ROM with 4 banks (64KB)
	rom := make([]uint8, 0x10000)
	// Bank 0
	rom[0x0100] = 0xB0
	// Bank 1
	rom[0x4100] = 0xB1
	// Bank 2
	rom[0x8100] = 0xB2
	// Bank 3
	rom[0xC100] = 0xB3

	// Create a RAM with 1 bank (8KB)
	ram := make([]uint8, 0x2000)

	mbc := CreateMBC1(rom, ram)

	// Test Initial State (Bank 0 at 0x0000-0x3FFF, Bank 1 at 0x4000-0x7FFF)
	if mbc.Read(0x0100) != 0xB0 {
		t.Errorf("Initial Bank 0: Expected 0xB0, got %02X", mbc.Read(0x0100))
	}
	if mbc.Read(0x4100) != 0xB1 {
		t.Errorf("Initial Bank 1: Expected 0xB1, got %02X", mbc.Read(0x4100))
	}

	// Test ROM Banking
	mbc.Write(0x2000, 0x02) // Select Bank 2
	if mbc.Read(0x4100) != 0xB2 {
		t.Errorf("Bank 2: Expected 0xB2, got %02X", mbc.Read(0x4100))
	}

	mbc.Write(0x2000, 0x03) // Select Bank 3
	if mbc.Read(0x4100) != 0xB3 {
		t.Errorf("Bank 3: Expected 0xB3, got %02X", mbc.Read(0x4100))
	}

	// Test Bank 0 mapping (0x00 -> 0x01)
	mbc.Write(0x2000, 0x00)
	if mbc.Read(0x4100) != 0xB1 {
		t.Errorf("Bank 0->1: Expected 0xB1, got %02X", mbc.Read(0x4100))
	}

	// Test RAM Enabling
	if mbc.Read(0xA000) != 0xFF {
		t.Error("RAM should be disabled by default (returning 0xFF)")
	}
	mbc.Write(0x0000, 0x0A) // Enable RAM
	mbc.Write(0xA000, 0x55)
	if mbc.Read(0xA000) != 0x55 {
		t.Errorf("RAM Write/Read: Expected 0x55, got %02X", mbc.Read(0xA000))
	}

	mbc.Write(0x0000, 0x00) // Disable RAM
	if mbc.Read(0xA000) != 0xFF {
		t.Error("RAM should return 0xFF when disabled")
	}

	// Test ROM Upper bits (RAM Bank Register)
	mbc.Write(0x4000, 0x01) // RAM Bank 1 / ROM bits 5-6
	// In ROM Mode 0 (default), RAM Bank Register affects ROM Upper Bits
	// For 4 banks, it shouldn't affect much if numRomBanks is small
}
