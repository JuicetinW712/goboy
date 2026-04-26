package cpu

func (c *CPU) executeBlock0(opcode uint8) TCycles {
	switch opcode {
	case 0x00: // NOP
		return 4
	case 0x01: // LD BC, imm16
		return c.LD_R16_IMM16(R16_BC)
	case 0x02: // LD [BC], A
		return c.LD_R16MEM_A(R16Mem_BC)
	case 0x03: // INC BC
		return c.INC_R16(R16_BC)
	case 0x04: // INC B
		return c.INC_R8(R8_B)
	case 0x05: // DEC B
		return c.DEC_R8(R8_B)
	case 0x06: // LD B, n8
		return c.LD_R8_IMM8(R8_B)
	case 0x07: // RLCA
		return c.RLCA()
	case 0x08: // LD [imm16], SP
		return c.LD_IMM16MEM_SP()
	case 0x09: // ADD HL, BC
		return c.ADD_HL_R16(R16_BC)
	case 0x0A: // LD A, [BC]
		return c.LD_A_R16MEM(R16Mem_BC)
	case 0x0B: // DEC BC
		return c.DEC_R16(R16_BC)
	case 0x0C: // INC C
		return c.INC_R8(R8_C)
	case 0x0D: // DEC C
		return c.DEC_R8(R8_C)
	case 0x0E: // LD C, imm8
		return c.LD_R8_IMM8(R8_C)
	case 0x0F: // RRCA
		return c.RRCA()
	case 0x10: // STOP
		return c.STOP()
	case 0x11: // LD DE, n16
		return c.LD_R16_IMM16(R16_DE)
	case 0x12: // LD [DE], A
		return c.LD_R16MEM_A(R16Mem_DE)
	case 0x13: // INC DE
		return c.INC_R16(R16_DE)
	case 0x14: // INC D
		return c.INC_R8(R8_D)
	case 0x15: // DEC D
		return c.DEC_R8(R8_D)
	case 0x16: // LD D, imm8
		return c.LD_R8_IMM8(R8_D)
	case 0x17: // RLA
		return c.RLA()
	case 0x18: // JR e8
		return c.JR_E8()
	case 0x19: // ADD HL, DE
		return c.ADD_HL_R16(R16_DE)
	case 0x1A: // LD A, [DE]
		return c.LD_A_R16MEM(R16Mem_DE)
	case 0x1B: // DEC DE
		return c.DEC_R16(R16_DE)
	case 0x1C: // INC E
		return c.INC_R8(R8_E)
	case 0x1D: // DEC E
		return c.DEC_R8(R8_E)
	case 0x1E: // LD E, imm8
		return c.LD_R8_IMM8(R8_E)
	case 0x1F: // RRA
		return c.RRA()
	case 0x20: // JR NZ, e8
		return c.JR_COND_E8(COND_NZ)
	case 0x21: // LD HL, n16
		return c.LD_R16_IMM16(R16_HL)
	case 0x22: // LD [HL+], A
		return c.LD_R16MEM_A(R16Mem_HLI)
	case 0x23: // INC HL
		return c.INC_R16(R16_HL)
	case 0x24: // INC H
		return c.INC_R8(R8_H)
	case 0x25: // DEC H
		return c.DEC_R8(R8_H)
	case 0x26: // LD H, n8
		return c.LD_R8_IMM8(R8_H)
	case 0x27: // DAA
		return c.DAA()
	case 0x28: // JR Z, e8
		return c.JR_COND_E8(COND_Z)
	case 0x29: // ADD HL, HL
		return c.ADD_HL_R16(R16_HL)
	case 0x2A: // LD A, [HL+]
		return c.LD_A_R16MEM(R16Mem_HLI)
	case 0x2B: // DEC HL
		return c.DEC_R16(R16_HL)
	case 0x2C: // INC L
		return c.INC_R8(R8_L)
	case 0x2D: // DEC L
		return c.DEC_R8(R8_L)
	case 0x2E: // LD L, n8
		return c.LD_R8_IMM8(R8_L)
	case 0x2F: // CPL
		return c.CPL()
	case 0x30: // JR NC, e8
		return c.JR_COND_E8(COND_NC)
	case 0x31: // LD SP, n16
		return c.LD_R16_IMM16(R16_SP)
	case 0x32: // LD [HL-], A
		return c.LD_R16MEM_A(R16Mem_HLD)
	case 0x33: // INC SP
		return c.INC_R16(R16_SP)
	case 0x34: // INC [HL]
		return c.INC_R8(R8_HL)
	case 0x35: // DEC [HL]
		return c.DEC_R8(R8_HL)
	case 0x36: // LD [HL], n8
		return c.LD_R8_IMM8(R8_HL)
	case 0x37: // SCF
		return c.SCF()
	case 0x38: // JR C, e8
		return c.JR_COND_E8(COND_C)
	case 0x39: // ADD HL, SP
		return c.ADD_HL_R16(R16_SP)
	case 0x3A: // LD A, [HL-]
		return c.LD_A_R16MEM(R16Mem_HLD)
	case 0x3B: // DEC SP
		return c.DEC_R16(R16_SP)
	case 0x3C: // INC A
		return c.INC_R8(R8_A)
	case 0x3D: // DEC A
		return c.DEC_R8(R8_A)
	case 0x3E: // LD A, n8
		return c.LD_R8_IMM8(R8_A)
	case 0x3F: // CCF
		return c.CCF()
	default:
		panic("Opcode should not have been assigned to block 0")
	}
}

func (c *CPU) LD_R16_IMM16(reg R16) TCycles {
	c.writeR16(reg, c.loadImm16())
	return 12
}

func (c *CPU) LD_R16MEM_A(reg R16Mem) TCycles {
	c.writeR16Mem(reg, c.readR8(R8_A))
	return 8
}

func (c *CPU) INC_R16(reg R16) TCycles {
	val := c.readR16(reg)
	c.writeR16(reg, val+1)
	return 8
}

func (c *CPU) INC_R8(reg R8) TCycles {
	val := c.readR8(reg)
	result := val + 1
	c.writeR8(reg, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag((val & 0x0F) == 0x0F)

	if reg == R8_HL {
		return 12
	} else {
		return 4
	}
}

func (c *CPU) DEC_R8(reg R8) TCycles {
	val := c.readR8(reg)
	result := val - 1
	c.writeR8(reg, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(true)
	c.setHalfCarryFlag((val & 0x0F) == 0x00)

	if reg == R8_HL {
		return 12
	} else {
		return 4
	}
}

func (c *CPU) LD_R8_IMM8(reg R8) TCycles {
	c.writeR8(reg, c.loadImm8())

	if reg == R8_HL {
		return 12
	} else {
		return 8
	}
}

func (c *CPU) LD_IMM16MEM_SP() TCycles {
	addr := c.loadImm16()
	low := uint8(c.SP & 0xFF)
	high := uint8(c.SP >> 8)

	c.bus.Write(addr, low)
	c.bus.Write(addr+1, high)
	return 20
}

func (c *CPU) ADD_HL_R16(reg R16) TCycles {
	hl := uint32(c.readR16(R16_HL))
	val := uint32(c.readR16(reg))

	result := hl + val
	c.writeR16(R16_HL, uint16(result))

	c.setSubtractionFlag(false)
	c.setHalfCarryFlag((hl&0xFFF)+(val&0xFFF) > 0xFFF)
	c.setCarryFlag(hl+val > 0xFFFF)
	return 8
}

func (c *CPU) LD_A_R16MEM(reg R16Mem) TCycles {
	val := c.readR16Mem(reg)
	c.writeR8(R8_A, val)
	return 8
}

func (c *CPU) DEC_R16(reg R16) TCycles {
	val := c.readR16(reg)
	c.writeR16(reg, val-1)
	return 8
}

func (c *CPU) RLCA() TCycles {
	value := c.readR8(R8_A)
	msb := value >> 7

	result := (value << 1) | msb
	c.writeR8(R8_A, result)

	c.setZeroFlag(false)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(msb == 1)
	return 4
}

func (c *CPU) RRCA() TCycles {
	value := c.readR8(R8_A)
	lsb := value & 0b1

	result := (value >> 1) | lsb<<7
	c.writeR8(R8_A, result)

	c.setZeroFlag(false)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(lsb == 1)
	return 4
}

func (c *CPU) RLA() TCycles {
	value := c.readR8(R8_A)
	msb := value >> 7

	var carry uint8
	if c.getCarryFlag() {
		carry = 1
	}

	result := (value << 1) | carry
	c.writeR8(R8_A, result)

	c.setZeroFlag(false)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(msb == 1)
	return 4
}

func (c *CPU) RRA() TCycles {
	value := c.readR8(R8_A)
	lsb := value & 0b1

	var carry uint8
	if c.getCarryFlag() {
		carry = 1
	}

	result := (value >> 1) | (carry << 7)
	c.writeR8(R8_A, result)

	c.setZeroFlag(false)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(lsb == 1)
	return 4
}

func (c *CPU) JR_E8() TCycles {
	offset := c.loadE8()
	c.PC = uint16(int32(c.PC) + int32(offset))
	return 12
}

func (c *CPU) STOP() TCycles {
	// Bug on og gameboy hardware where the STOP instruction causes
	// the CPU to skip the next byte in ROM
	c.loadImm8()
	c.Stopped = true
	return 4
}

func (c *CPU) JR_COND_E8(condition COND) TCycles {
	offset := c.loadE8()

	if c.evaluateCondition(condition) {
		c.PC = uint16(int32(c.PC) + int32(offset))
		return 12
	}

	return 8
}

func (cpu *CPU) DAA() TCycles {
	a := cpu.readR8(R8_A)
	correction := uint8(0)
	c := cpu.getCarryFlag()

	if cpu.getSubtractionFlag() {
		if cpu.getHalfCarryFlag() {
			correction |= 0x06
		}

		if c {
			correction |= 0x60
		}

		a -= correction
	} else {
		if cpu.getHalfCarryFlag() || (a&0x0F) > 0x09 {
			correction |= 0x06
		}

		if c || a > 0x99 {
			correction |= 0x60
			c = true
		}

		a += correction
	}

	cpu.setZeroFlag(a == 0)
	cpu.setHalfCarryFlag(false)
	cpu.setCarryFlag(c)
	cpu.writeR8(R8_A, a)
	return 4
}

func (c *CPU) CPL() TCycles {
	value := c.readR8(R8_A)
	c.writeR8(R8_A, ^value)

	c.setSubtractionFlag(true)
	c.setHalfCarryFlag(true)
	return 4
}

func (c *CPU) SCF() TCycles {
	c.setCarryFlag(true)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	return 4
}

func (c *CPU) CCF() TCycles {
	c.setCarryFlag(!c.getCarryFlag())
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	return 4
}
