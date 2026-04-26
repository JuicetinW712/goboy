package cpu

func (c *CPU) executeBlockCB(opcode uint8) TCycles {
	registerIdx := opcode & 0b111
	targetRegister := c.getR8FromIdx(registerIdx)

	bitIdx := (opcode >> 3) & 0b111
	switch opcode >> 6 {
	case 0: // Rotates and shifts
		switch opcode >> 3 {
		case 0:
			return c.RLC(targetRegister)
		case 1:
			return c.RRC(targetRegister)
		case 2:
			return c.RL(targetRegister)
		case 3:
			return c.RR(targetRegister)
		case 4:
			return c.SLA(targetRegister)
		case 5:
			return c.SRA(targetRegister)
		case 6:
			return c.SWAP(targetRegister)
		case 7:
			return c.SRL(targetRegister)
		default:
			panic("Invalid Rotate/Shift opcode")
		}
	case 1:
		return c.BIT(bitIdx, targetRegister)
	case 2:
		return c.RES(bitIdx, targetRegister)
	case 3:
		return c.SET(bitIdx, targetRegister)
	default:
		panic("Invalid opcode in block CB")
	}
}

func (c *CPU) RLC(targetRegister R8) TCycles {
	value := c.readR8(targetRegister)
	msb := value >> 7

	result := (value << 1) | msb
	c.writeR8(targetRegister, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(msb == 1)

	if targetRegister == R8_HL {
		return 16
	}
	return 8
}

func (c *CPU) RRC(targetRegister R8) TCycles {
	value := c.readR8(targetRegister)
	lsb := value & 0x1

	result := (value >> 1) | (lsb << 7)
	c.writeR8(targetRegister, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(lsb == 1)

	if targetRegister == R8_HL {
		return 16
	}
	return 8
}

func (c *CPU) RL(targetRegister R8) TCycles {
	value := c.readR8(targetRegister)
	msb := value >> 7

	var carry uint8
	if c.getCarryFlag() {
		carry = 1
	}

	result := (value << 1) | carry
	c.writeR8(targetRegister, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(msb == 1)

	if targetRegister == R8_HL {
		return 16
	}
	return 8
}

func (c *CPU) RR(targetRegister R8) TCycles {
	value := c.readR8(targetRegister)
	lsb := value & 0x1

	var carry uint8
	if c.getCarryFlag() {
		carry = 1
	}

	result := (value >> 1) | (carry << 7)
	c.writeR8(targetRegister, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(lsb == 1)

	if targetRegister == R8_HL {
		return 16
	}
	return 8
}

func (c *CPU) SLA(targetRegister R8) TCycles {
	value := c.readR8(targetRegister)
	msb := value >> 7

	result := value << 1
	c.writeR8(targetRegister, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(msb == 1)

	if targetRegister == R8_HL {
		return 16
	}
	return 8
}

func (c *CPU) SRA(targetRegister R8) TCycles {
	value := c.readR8(targetRegister)
	lsb := value & 0x1

	result := (value >> 1) | (value & 0x80)
	c.writeR8(targetRegister, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(lsb == 1)

	if targetRegister == R8_HL {
		return 16
	}
	return 8
}

func (c *CPU) SWAP(targetRegister R8) TCycles {
	value := c.readR8(targetRegister)

	result := (value << 4) | (value >> 4)
	c.writeR8(targetRegister, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(false)

	if targetRegister == R8_HL {
		return 16
	}
	return 8
}

func (c *CPU) SRL(targetRegister R8) TCycles {
	value := c.readR8(targetRegister)
	lsb := value & 0x1

	result := value >> 1
	c.writeR8(targetRegister, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(lsb == 1)

	if targetRegister == R8_HL {
		return 16
	}
	return 8
}

func (c *CPU) BIT(bitIdx uint8, targetRegister R8) TCycles {
	value := c.readR8(targetRegister)
	c.setZeroFlag((value>>bitIdx)&0x1 == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(true)

	if targetRegister == R8_HL {
		return 12
	}
	return 8
}

func (c *CPU) RES(bitIdx uint8, targetRegister R8) TCycles {
	value := c.readR8(targetRegister)
	result := value & ^(1 << bitIdx)
	c.writeR8(targetRegister, result)

	if targetRegister == R8_HL {
		return 16
	}
	return 8
}

func (c *CPU) SET(bitIdx uint8, reg R8) TCycles {
	value := c.readR8(reg)
	result := value | (1 << bitIdx)
	c.writeR8(reg, result)

	if reg == R8_HL {
		return 16
	}
	return 8
}
