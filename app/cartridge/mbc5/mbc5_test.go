package mbc5

import (
	"testing"
)

func TestMBC5(t *testing.T) {
	rom := make([]uint8, 0x10000) // 4 banks
	rom[0x0100] = 0xB0
	rom[0x4100] = 0xB1
	rom[0x8100] = 0xB2
	
	ram := make([]uint8, 0x4000) // 2 banks
	
	mbc := CreateMBC5(rom, ram)
	
	// Initial state
	if mbc.Read(0x0100) != 0xB0 {
		t.Errorf("Expected 0xB0, got %02X", mbc.Read(0x0100))
	}
	
	// ROM Bank 1
	mbc.Write(0x2000, 0x01)
	if mbc.Read(0x4100) != 0xB1 {
		t.Errorf("Expected 0xB1, got %02X", mbc.Read(0x4100))
	}
	
	// ROM Bank 2
	mbc.Write(0x2000, 0x02)
	if mbc.Read(0x4100) != 0xB2 {
		t.Errorf("Expected 0xB2, got %02X", mbc.Read(0x4100))
	}
	
	// ROM Bank 9-th bit
	mbc.Write(0x3000, 0x01) // Set bit 8
	// numRomBanks = 4, so bank register will be 2 (from previous) | 256 = 258
	// 258 % 4 = 2.
	if mbc.romBankRegister != 258 {
		t.Errorf("Expected 258, got %d", mbc.romBankRegister)
	}

	// RAM Enable
	mbc.Write(0x0000, 0x0A)
	mbc.Write(0x4000, 0x01) // RAM Bank 1
	mbc.Write(0xA100, 0x55)
	if mbc.Read(0xA100) != 0x55 {
		t.Errorf("Expected 0x55, got %02X", mbc.Read(0xA100))
	}
}
