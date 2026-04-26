package ppu

// 6 - LYC int select (Read/Write): If set, selects the LYC == LY condition for the STAT interrupt.
// 5 - Mode 2 int select (Read/Write): If set, selects the Mode 2 condition for the STAT interrupt.
// 4 - Mode 1 int select (Read/Write): If set, selects the Mode 1 condition for the STAT interrupt.
// 3 - Mode 0 int select (Read/Write): If set, selects the Mode 0 condition for the STAT interrupt.
// 2 - LYC == LY (Read-only): Set when LY contains the same value as LYC; it is constantly updated.
// 1,0 - PPU mode (Read-only): Indicates the PPU’s current status. Reports 0 instead when the PPU is disabled.

// Need to emulate STAT blocking behavior
// STAT is essentially a big or statement that only
// triggers if moving from low to high
//
// STAT_line =
//     (Mode0 && enabled)
//  || (Mode1 && enabled)
//  || (Mode2 && enabled)
//  || (LY == LYC && enabled)

func (p *PPU) statInterruptLine() bool {
	if !p.lcdEnabled() {
		return false
	}

	mode := p.getMode()

	// Check if any enabled interrupt source is currently "high"
	lycBit := (p.STAT>>6)&1 == 1 && (p.LY == p.LYC)
	m2Bit := (p.STAT>>5)&1 == 1 && mode == OAMScan
	m1Bit := (p.STAT>>4)&1 == 1 && mode == VBlank
	m0Bit := (p.STAT>>3)&1 == 1 && mode == HBlank

	return lycBit || m2Bit || m1Bit || m0Bit
}

func (p *PPU) checkStatInterrupt() {
	current := p.statInterruptLine()

	if !p.lowStatValue && current {
		p.OnStat()
	}

	p.lowStatValue = current
}

func (p *PPU) updateLYC() {
	if p.LY == p.LYC {
		p.STAT = p.STAT | 0b00000100
	} else {
		p.STAT = p.STAT & 0b11111011
	}
}

func (p *PPU) setMode(mode PPU_Mode) {
	p.STAT = (p.STAT & 0b11111100) | uint8(mode)
}

func (p *PPU) getMode() PPU_Mode {
	return PPU_Mode(p.STAT & 0b11)
}
