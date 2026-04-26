package joypad

import (
	"testing"
)

func TestJoypad(t *testing.T) {
	j := CreateJoypad()
	
	// Initial state: j.sel = 0x00, j.dpad = 0x0F, j.buttons = 0x0F
	// j.Read() returns 0xC0 | j.sel | (res)
	// res = 0x0F & dpad & buttons = 0x0F
	// return 0xC0 | 0x00 | 0x0F = 0xCF
	if j.Read() != 0xCF {
		t.Errorf("Initial read failed, got %02X", j.Read())
	}
	
	// Press A (Button 4)
	j.PressButton(BUTTON_A)
	// j.buttons = 0x0F &^ (1 << (4-4)) = 0x0E
	// Select buttons (bit 5 = 0, bit 4 = 1)
	j.Write(0x10) 
	// res = 0x0F & j.buttons = 0x0E
	// return 0xC0 | 0x10 | 0x0E = 0xDE
	if j.Read() != 0xDE {
		t.Errorf("Press A failed, got %02X", j.Read())
	}
	
	// Select dpad (bit 5 = 1, bit 4 = 0)
	j.Write(0x20)
	// res = 0x0F & j.dpad = 0x0F
	// return 0xC0 | 0x20 | 0x0F = 0xEF
	if j.Read() != 0xEF {
		t.Errorf("Select dpad failed, got %02X", j.Read())
	}
	
	// Press Left (Button 1)
	j.PressButton(DPAD_LEFT)
	// j.dpad = 0x0F &^ (1 << 1) = 0x0D
	// res = 0x0F & j.dpad = 0x0D
	// return 0xC0 | 0x20 | 0x0D = 0xED
	if j.Read() != 0xED {
		t.Errorf("Press Left failed, got %02X", j.Read())
	}
	
	// Release Left
	j.ReleaseButton(DPAD_LEFT)
	if j.Read() != 0xEF {
		t.Errorf("Release Left failed, got %02X", j.Read())
	}
	
	// Release A
	j.ReleaseButton(BUTTON_A)
	j.Write(0x10)
	if j.Read() != 0xDF {
		t.Errorf("Release A failed, got %02X", j.Read())
	}
}
