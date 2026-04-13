package cartridge

import (
	"testing"
)

func TestCartridgeInfo(t *testing.T) {
	rom := make([]uint8, 0x8000)
	copy(rom[TITLE_START_ADDR:TITLE_END_ADDR], "TESTROM")
	rom[MBC_TYPE_ADDR] = 0x00 // ROM ONLY
	rom[RAM_SIZE_ADDR] = 0x00 // No RAM

	cart := NewCartridge("test.gb", rom)
	if cart.title != "TESTROM" {
		t.Errorf("Expected title TESTROM, got %s", cart.title)
	}

	info := cart.String()
	if info == "" {
		t.Error("String() returned empty info")
	}
}

func TestNewCartridgeFromFile(t *testing.T) {
	// This would require a real file, skipping for now or mocking os.ReadFile
}
