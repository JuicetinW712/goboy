package bus

import (
	"testing"
)

type MockDevice struct {
	data map[uint16]uint8
}

func NewMockDevice() *MockDevice {
	return &MockDevice{data: make(map[uint16]uint8)}
}

func (m *MockDevice) Read(addr uint16) uint8 {
	return m.data[addr]
}

func (m *MockDevice) Write(addr uint16, val uint8) {
	m.data[addr] = val
}

func TestMMUReadWrite(t *testing.T) {
	cart := NewMockDevice()
	ppu := NewMockDevice()
	io := NewMockDevice()
	mmu := NewMMU(cart, ppu, io)

	tests := []struct {
		name     string
		addr     uint16
		val      uint8
		device   *MockDevice
		readAddr uint16 // some regions map to different device addrs
	}{
		{"ROM 0", 0x0000, 0x11, cart, 0x0000},
		{"ROM N", 0x4000, 0x22, cart, 0x4000},
		{"VRAM", 0x8000, 0x33, ppu, 0x8000},
		{"ERAM", 0xA000, 0x44, cart, 0xA000},
		{"WRAM", 0xC000, 0x55, nil, 0}, // Internal to MMU
		{"Echo RAM", 0xE000, 0x66, nil, 0}, // Internal to MMU (maps to WRAM)
		{"OAM", 0xFE00, 0x77, ppu, 0xFE00},
		{"IO", 0xFF00, 0x88, io, 0xFF00},
		{"HRAM", 0xFF80, 0x99, nil, 0}, // Internal to MMU
		{"IE", 0xFFFF, 0xAA, nil, 0}, // Internal to MMU
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmu.Write(tt.addr, tt.val)
			got := mmu.Read(tt.addr)
			if got != tt.val {
				t.Errorf("Read/Write failed at 0x%04X: expected 0x%02X, got 0x%02X", tt.addr, tt.val, got)
			}

			if tt.device != nil {
				if tt.device.Read(tt.readAddr) != tt.val {
					t.Errorf("Device Read failed at 0x%04X: expected 0x%02X, got 0x%02X", tt.readAddr, tt.val, tt.device.Read(tt.readAddr))
				}
			}
		})
	}
}

func TestEchoRAM(t *testing.T) {
	mmu := NewMMU(NewMockDevice(), NewMockDevice(), NewMockDevice())
	
	// Write to WRAM, read from Echo RAM
	mmu.Write(0xC000, 0xAB)
	if mmu.Read(0xE000) != 0xAB {
		t.Errorf("Echo RAM read failed: expected 0xAB, got %02X", mmu.Read(0xE000))
	}

	// Write to Echo RAM, read from WRAM
	mmu.Write(0xFDFF, 0xCD)
	if mmu.Read(0xDDFF) != 0xCD {
		t.Errorf("WRAM read from Echo RAM write failed: expected 0xCD, got %02X", mmu.Read(0xDDFF))
	}
}

func TestUnusableMemory(t *testing.T) {
	mmu := NewMMU(NewMockDevice(), NewMockDevice(), NewMockDevice())
	
	// Should return 0xFF and ignore writes
	addr := uint16(0xFEA0)
	mmu.Write(addr, 0x55)
	if mmu.Read(addr) != 0xFF {
		t.Errorf("Unusable memory should return 0xFF, got %02X", mmu.Read(addr))
	}
}
