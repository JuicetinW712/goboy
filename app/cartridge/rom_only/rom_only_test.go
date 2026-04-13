package rom_only

import (
	"testing"
)

func TestROMOnlyMBC(t *testing.T) {
	rom := make([]uint8, 0x8000)
	rom[0x1234] = 0x55
	mbc := CreateROMOnlyMBC(rom)

	if mbc.Read(0x1234) != 0x55 {
		t.Errorf("Expected 0x55, got %02X", mbc.Read(0x1234))
	}

	mbc.Write(0x1234, 0xAA)
	if mbc.Read(0x1234) != 0x55 {
		t.Error("ROM Only MBC should not be writeable")
	}
}
