package ppu

import (
	"testing"
)

func TestPPUModeTransitions(t *testing.T) {
	ppu := CreatePPU()
	vblankCount := 0
	ppu.OnVBlank = func() {
		vblankCount++
	}
	ppu.OnStat = func() {}
	
	// Initial state: LCD on, Mode HBlank
	ppu.LCDC = 0x80
	
	// Dots 0-79: OAM Scan
	ppu.Step(10)
	if ppu.getMode() != OAMScan {
		t.Errorf("Expected OAMScan, got %d", ppu.getMode())
	}
	
	// Dots 80-251: Drawing
	ppu.Step(70) // Total dots 80
	if ppu.getMode() != Drawing {
		t.Errorf("Expected Drawing, got %d", ppu.getMode())
	}
	
	// Dots 252-455: HBlank
	ppu.Step(172) // Total dots 252
	if ppu.getMode() != HBlank {
		t.Errorf("Expected HBlank, got %d", ppu.getMode())
	}
	
	// Dot 456: Next Line (LY=1)
	ppu.Step(204) // Total dots 456
	if ppu.LY != 1 {
		t.Errorf("Expected LY=1, got %d", ppu.LY)
	}
	
	// Reach VBlank (LY=144)
	for i := 1; i < 144; i++ {
		ppu.Step(456)
	}
	
	if ppu.LY != 144 {
		t.Errorf("Expected LY=144, got %d", ppu.LY)
	}
	if ppu.getMode() != VBlank {
		t.Errorf("Expected VBlank, got %d", ppu.getMode())
	}
	if vblankCount != 1 {
		t.Errorf("Expected vblankCount=1, got %d", vblankCount)
	}
}

func TestPPUReadWrite(t *testing.T) {
	ppu := CreatePPU()
	
	// Write VRAM
	ppu.Write(0x8000, 0x42)
	if ppu.Read(0x8000) != 0x42 {
		t.Errorf("VRAM Read/Write failed, got %02X", ppu.Read(0x8000))
	}

	// Write OAM
	ppu.Write(0xFE00, 0x55)
	if ppu.Read(0xFE00) != 0x55 {
		t.Errorf("OAM Read/Write failed, got %02X", ppu.Read(0xFE00))
	}

	// Mode 3 (Drawing) should block VRAM and OAM access
	ppu.LCDC = 0x80
	ppu.Step(100) // Mode Drawing (80 to 251 dots)
	
	if ppu.Read(0x8000) != 0xFF {
		t.Errorf("VRAM should be blocked during Drawing mode, got %02X", ppu.Read(0x8000))
	}
	if ppu.Read(0xFE00) != 0xFF {
		t.Errorf("OAM should be blocked during Drawing mode, got %02X", ppu.Read(0xFE00))
	}

	// Test STAT write
	ppu.Write(STAT_ADDRESS, 0xF8) // Set all interrupt enable bits
	if (ppu.STAT & 0xF8) != 0xF8 {
		t.Errorf("STAT write failed, got %02X", ppu.STAT)
	}
}

func TestPPURendering(t *testing.T) {
	ppu := CreatePPU()
	ppu.LCDC = 0x80 | 0x10 | 0x01 // LCD on, Tile Data 0x8000, BG enabled
	ppu.BGP = 0xE4 // 11 10 01 00
	
	// Set up tile 0 at 0x8000
	ppu.VRAM[0] = 0x80
	ppu.VRAM[1] = 0x00
	
	// Set up tile 1 at 0x8010
	ppu.VRAM[0x10] = 0x00
	ppu.VRAM[0x11] = 0x80

	// Tile map at 0x9800
	ppu.VRAM[0x9800-0x8000] = 0x00 // Tile 0 at (0,0)
	ppu.VRAM[0x9801-0x8000] = 0x01 // Tile 1 at (1,0)
	
	ppu.LY = 0
	ppu.drawBackgroundLine()
	
	c0 := ppu.Buffer.RGBAAt(0, 0)
	if c0 != colors[1] {
		t.Errorf("Expected colors[1], got %v", c0)
	}
	
	c1 := ppu.Buffer.RGBAAt(8, 0)
	if c1 != colors[2] {
		t.Errorf("Expected colors[2], got %v", c1)
	}

	// Test Window
	ppu.LCDC |= 0x20 // Window enable
	ppu.WY = 0
	ppu.WX = 7
	ppu.VRAM[0x9800-0x8000] = 0x01 // Tile 1 at (0,0) in window map
	ppu.drawWindowLine()
	cW := ppu.Buffer.RGBAAt(0, 0)
	if cW != colors[2] {
		t.Errorf("Expected colors[2] for window, got %v", cW)
	}

	// Test Sprites
	ppu.LCDC |= 0x02 // OBJ enable
	ppu.OBP0 = 0xE4
	// Sprite 0: Y=16 (LY=0), X=8 (Screen X=0), Tile=0, Attr=0
	ppu.OAM[0] = 16
	ppu.OAM[1] = 8
	ppu.OAM[2] = 0
	ppu.OAM[3] = 0
	ppu.drawSpriteLine()
	cS := ppu.Buffer.RGBAAt(0, 0)
	if cS != colors[1] {
		t.Errorf("Expected colors[1] for sprite, got %v", cS)
	}

	// Test Scrolling
	ppu.SCX = 8
	ppu.SCY = 0
	ppu.drawBackgroundLine() // (0,0) on screen should now be (8,0) in tile map -> Tile 1
	cSc := ppu.Buffer.RGBAAt(0, 0)
	if cSc != colors[2] {
		t.Errorf("Expected colors[2] for scrolled BG, got %v", cSc)
	}
}

func TestPPUStatInterrupt(t *testing.T) {
	ppu := CreatePPU()
	ppu.LCDC = 0x80
	ppu.STAT = 0x00
	
	statFired := false
	ppu.OnStat = func() {
		statFired = true
	}

	// Enable LYC=LY interrupt
	ppu.Write(STAT_ADDRESS, 0x40)
	ppu.LYC = 10
	ppu.LY = 10
	ppu.updateLYC()
	ppu.checkStatInterrupt()
	
	if !statFired {
		t.Error("STAT interrupt should have fired for LYC=LY")
	}
	
	statFired = false
	ppu.checkStatInterrupt()
	if statFired {
		t.Error("STAT interrupt should not fire again (no transition)")
	}

	ppu.LY = 11
	ppu.updateLYC()
	ppu.checkStatInterrupt()
	statFired = false
	
	ppu.LY = 10
	ppu.updateLYC()
	ppu.checkStatInterrupt()
	if !statFired {
		t.Error("STAT interrupt should have fired again for LYC=LY transition")
	}
}

func TestPPULCDCMethods(t *testing.T) {
	ppu := CreatePPU()
	ppu.LCDC = 0xFF
	if !ppu.lcdEnabled() || !ppu.windowEnabled() || !ppu.objEnabled() || !ppu.backgroundAndWindowEnabled() {
		t.Error("LCDC methods failed")
	}
	if ppu.getWindowTileMap() != 1 || ppu.getBackgroundTileMap() != 1 || ppu.getObjSize() != 16 || !ppu.unsignedTileAddressing() {
		t.Errorf("LCDC complex methods failed: maps=[%d,%d] obj=%d unsig=%v", 
			ppu.getWindowTileMap(), ppu.getBackgroundTileMap(), ppu.getObjSize(), ppu.unsignedTileAddressing())
	}
	
	ppu.LCDC = 0x00
	if ppu.unsignedTileAddressing() {
		t.Error("unsignedTileAddressing should be false when bit 4 is 0")
	}
}

func TestPPUGetColor(t *testing.T) {
	ppu := CreatePPU()
	palette := uint8(0xE4) // 11 10 01 00
	if ppu.getColor(0, palette) != colors[0] { t.Error("Color 0 failed") }
	if ppu.getColor(1, palette) != colors[1] { t.Error("Color 1 failed") }
	if ppu.getColor(2, palette) != colors[2] { t.Error("Color 2 failed") }
	if ppu.getColor(3, palette) != colors[3] { t.Error("Color 3 failed") }
}
