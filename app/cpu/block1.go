package cpu

func (c *CPU) executeBlock1(opcode uint8) TCycles {
	// HALT
	if opcode == 0x76 {
		interruptFlags := c.bus.Read(0xFF0F)
		enabledInterrupts := c.bus.Read(0xFFFF)
		pendingInterrupts := interruptFlags & enabledInterrupts

		if !c.interruptsEnabled {
			if pendingInterrupts != 0 {
				// HALT bug: PC is not incremented for the next instruction
				c.haltBug = true
			} else {
				c.Halted = true
			}
		} else {
			if pendingInterrupts != 0 {
				// Interrupt pending and IME=1: CPU doesn't halt, services interrupt
				// In our implementation, Step will handle this as it checks interrupts before fetch.
				// But we should NOT set c.Halted here.
			} else {
				c.Halted = true
			}
		}
		return 4
	}

	// LD r8, r8
	dst := c.getR8FromIdx((opcode >> 3) & 0b111)
	src := c.getR8FromIdx(opcode & 0b111)

	srcValue := c.readR8(src)
	c.writeR8(dst, srcValue)

	if src == R8_HL || dst == R8_HL {
		return 8
	}
	return 4
}
