package cpu

func (c *CPU) executeBlock2(opcode uint8) TCycles {
	operationIdx := (opcode >> 3) & 0b111
	registerIdx := opcode & 0b111

	r8 := c.getR8FromIdx(registerIdx)
	value := c.readR8(r8)

	switch operationIdx {
	case 0x0:
		c.ADD(value)
	case 0x1:
		c.ADC(value)
	case 0x2:
		c.SUB(value)
	case 0x3:
		c.SBC(value)
	case 0x4:
		c.AND(value)
	case 0x5:
		c.XOR(value)
	case 0x6:
		c.OR(value)
	case 0x7:
		c.CP(value)
	default:
		panic("Invalid opcode in block 2")
	}

	return If(r8 == R8_HL, 8, 4)
}

func (c *CPU) ADD(value uint8) {
	a := c.readR8(R8_A)

	result := a + value
	c.writeR8(R8_A, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag((a&0xF)+(value&0xF) > 0xF)
	c.setCarryFlag(result < a)
}

func (c *CPU) ADC(value uint8) {
	a := c.readR8(R8_A)

	var carry uint8
	if c.getCarryFlag() {
		carry = 1
	}

	result := a + value + carry
	c.writeR8(R8_A, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag((a&0xF)+(value&0xF)+carry > 0xF)
	c.setCarryFlag(uint16(a)+uint16(value)+uint16(carry) > 0xFF)
}

func (c *CPU) SUB(value uint8) {
	a := c.readR8(R8_A)

	result := a - value
	c.writeR8(R8_A, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(true)
	c.setHalfCarryFlag((a & 0xF) < (value & 0xF))
	c.setCarryFlag(value > a)
}

func (c *CPU) SBC(value uint8) {
	a := c.readR8(R8_A)

	var carry uint8
	if c.getCarryFlag() {
		carry = 1
	}

	diff := uint16(a) - uint16(value) - uint16(carry)
	result := uint8(diff)

	c.writeR8(R8_A, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(true)
	c.setHalfCarryFlag((a & 0xF) < (value&0xF)+carry)
	c.setCarryFlag(uint16(a) < uint16(value)+uint16(carry))
}

func (c *CPU) AND(value uint8) {
	result := c.readR8(R8_A) & value
	c.writeR8(R8_A, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(true)
	c.setCarryFlag(false)
}

func (c *CPU) XOR(value uint8) {
	result := c.readR8(R8_A) ^ value
	c.writeR8(R8_A, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(false)
}

func (c *CPU) OR(value uint8) {
	result := c.readR8(R8_A) | value
	c.writeR8(R8_A, result)

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag(false)
	c.setCarryFlag(false)
}

func (c *CPU) CP(value uint8) {
	a := c.readR8(R8_A)

	result := a - value

	c.setZeroFlag(result == 0)
	c.setSubtractionFlag(true)
	c.setHalfCarryFlag((a & 0xF) < (value & 0xF))
	c.setCarryFlag(value > a)
}

func If(cond bool, vtrue, vfalse TCycles) TCycles {
	if cond {
		return vtrue
	}
	return vfalse
}
