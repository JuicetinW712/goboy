package mbc3

import (
	"testing"
)

func TestMBC3(t *testing.T) {
	rom := make([]uint8, 0x10000) // 4 banks
	rom[0x0100] = 0xB0
	rom[0x4100] = 0xB1
	rom[0x8100] = 0xB2
	
	ram := make([]uint8, 0x4000) // 2 banks
	
	mbc := CreateMBC3(rom, ram)
	
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
	
	// RAM Enable
	mbc.Write(0x0000, 0x0A)
	mbc.Write(0x4000, 0x01) // RAM Bank 1
	mbc.Write(0xA100, 0x55)
	if mbc.Read(0xA100) != 0x55 {
		t.Errorf("Expected 0x55, got %02X", mbc.Read(0xA100))
	}
	
	// RAM Disable
	mbc.Write(0x0000, 0x00)
	if mbc.Read(0xA100) != 0xFF {
		t.Errorf("Expected 0xFF, got %02X", mbc.Read(0xA100))
	}

	// Test RTC registers access (currently should return 0xFF)
	mbc.Write(0x0000, 0x0A) // Enable
	mbc.Write(0x4000, 0x08) // Select RTC S
	if mbc.Read(0xA000) != 0xFF {
		t.Errorf("Expected 0xFF for RTC, got %02X", mbc.Read(0xA000))
	}
	mbc.Write(0xA000, 0x12) // Write to RTC (should do nothing currently)

	// Latch clock (should do nothing currently)
	mbc.Write(0x6000, 0x00)
	mbc.Write(0x6000, 0x01)
}
