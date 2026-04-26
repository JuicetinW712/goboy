package cpu

func (c *CPU) executeBlock3(opcode uint8) TCycles {
	switch opcode {
	case 0xC0: // RET NZ
		return c.RET_COND(COND_NZ)
	case 0xC1: // POP BC
		return c.POP_R16STK(R16Stk_BC)
	case 0xC2: // JP NZ, a16
		return c.JP_COND_IMM16(COND_NZ)
	case 0xC3: // JP a16
		return c.JP_IMM16()
	case 0xC4: // CALL NZ, a16
		return c.CALL_COND_IMM16(COND_NZ)
	case 0xC5: // PUSH BC
		return c.PUSH_R16STK(R16Stk_BC)
	case 0xC6: // ADD A, n8
		c.ADD(c.loadImm8())
		return 8
	case 0xC7: // RST $00
		return c.RST_TGT3(0x00)
	case 0xC8: // RET Z
		return c.RET_COND(COND_Z)
	case 0xC9: // RET
		return c.RET()
	case 0xCA: // JP Z, a16
		return c.JP_COND_IMM16(COND_Z)
	case 0xCC: // CALL Z, a16
		return c.CALL_COND_IMM16(COND_Z)
	case 0xCD: // CALL a16
		return c.CALL_IMM16()
	case 0xCE: // ADC A, n8
		c.ADC(c.loadImm8())
		return 8
	case 0xCF: // RST $08
		return c.RST_TGT3(0x08)
	case 0xD0: // RET NC
		return c.RET_COND(COND_NC)
	case 0xD1: // POP DE
		return c.POP_R16STK(R16Stk_DE)
	case 0xD2: // JP NC, a16
		return c.JP_COND_IMM16(COND_NC)
	case 0xD4: // CALL NC, a16
		return c.CALL_COND_IMM16(COND_NC)
	case 0xD5: // PUSH DE
		return c.PUSH_R16STK(R16Stk_DE)
	case 0xD6: // SUB A, n8
		c.SUB(c.loadImm8())
		return 8
	case 0xD7: // RST $10
		return c.RST_TGT3(0x10)
	case 0xD8: // RET C
		return c.RET_COND(COND_C)
	case 0xD9: // RETI
		return c.RETI()
	case 0xDA: // JP C, a16
		return c.JP_COND_IMM16(COND_C)
	case 0xDC: // CALL C, imm16
		return c.CALL_COND_IMM16(COND_C)
	case 0xDE: // SBC A, n8
		c.SBC(c.loadImm8())
		return 8
	case 0xDF: // RST $18
		return c.RST_TGT3(0x18)
	case 0xE0: // LD (FF00 + imm8), A
		return c.LDH_IMM8MEM_A()
	case 0xE1: // POP HL
		return c.POP_R16STK(R16Stk_HL)
	case 0xE2: // LD (FF00 + C), A
		return c.LDH_CMEM_A()
	case 0xE5: // PUSH HL
		return c.PUSH_R16STK(R16Stk_HL)
	case 0xE6: // AND A, n8
		c.AND(c.loadImm8())
		return 8
	case 0xE7: // RST $20
		return c.RST_TGT3(0x20)
	case 0xE8: // ADD SP, e8
		return c.ADD_SP_E8()
	case 0xE9: // JP HL
		return c.JP_HL()
	case 0xEA: // LD [imm16], A
		return c.LD_IMM16MEM_A()
	case 0xEE: // XOR A, imm8
		c.XOR(c.loadImm8())
		return 8
	case 0xEF: // RST $28
		return c.RST_TGT3(0x28)
	case 0xF0: // LD A, (FF00 + imm8)
		return c.LDH_A_IMM8MEM()
	case 0xF1: // POP AF
		return c.POP_R16STK(R16Stk_AF)
	case 0xF2: // LD A, (FF00 + C)
		return c.LDH_A_C()
	case 0xF3: // DI
		return c.DI()
	case 0xF5: // PUSH AF
		return c.PUSH_R16STK(R16Stk_AF)
	case 0xF6: // OR A, imm8
		c.OR(c.loadImm8())
		return 8
	case 0xF7: // RST $30
		return c.RST_TGT3(0x30)
	case 0xF8: // LD HL, SP + e8
		return c.LD_HL_SP_PLUS_E8()
	case 0xF9: // LD SP, HL
		return c.LD_SP_HL()
	case 0xFA: // LD A, [imm16]
		return c.LD_A_IMM16MEM()
	case 0xFB: // EI
		return c.EI()
	case 0xFE: // CP A, imm8
		c.CP(c.loadImm8())
		return 8
	case 0xFF:
		return c.RST_TGT3(0x38)
	default:
		panic("Invalid opcode in block 3")
	}
}

func (c *CPU) RET() TCycles {
	c.PC = c.popStack16()
	return 16
}

func (c *CPU) RET_COND(condition COND) TCycles {
	if c.evaluateCondition(condition) {
		c.RET()
		return 20
	}
	return 8
}

func (c *CPU) POP_R16STK(r R16Stk) TCycles {
	value := c.popStack16()
	c.writeR16Stk(r, value)
	return 12
}

func (c *CPU) JP_COND_IMM16(condition COND) TCycles {
	value := c.loadImm16()

	if c.evaluateCondition(condition) {
		c.PC = value
		return 16
	}
	return 12
}

func (c *CPU) JP_IMM16() TCycles {
	value := c.loadImm16()
	c.PC = value
	return 16
}

func (c *CPU) CALL_COND_IMM16(condition COND) TCycles {
	address := c.loadImm16() // Advance the PC by 2 bytes regardless!

	if c.evaluateCondition(condition) {
		c.call(address)
		return 24
	}
	return 12
}

func (c *CPU) CALL_IMM16() TCycles {
	address := c.loadImm16()
	c.call(address)
	return 24
}

func (c *CPU) call(address uint16) {
	c.pushStack16(c.PC)
	c.PC = address
}

func (c *CPU) PUSH_R16STK(r R16Stk) TCycles {
	value := c.readR16Stk(r)
	c.pushStack16(value)
	return 16
}

func (c *CPU) RST_TGT3(address uint16) TCycles {
	c.pushStack16(c.PC)
	c.PC = address
	return 16
}

func (c *CPU) RETI() TCycles {
	cycles := c.RET()
	c.interruptsEnabled = true
	return cycles
}

func (c *CPU) LDH_IMM8MEM_A() TCycles {
	offset := c.loadImm8()
	value := c.readR8(R8_A)
	c.bus.Write(uint16(0xFF00)+uint16(offset), value)
	return 12
}

func (c *CPU) LDH_A_IMM8MEM() TCycles {
	offset := c.loadImm8()
	value := c.bus.Read(0xFF00 + uint16(offset))
	c.writeR8(R8_A, value)
	return 12
}

func (c *CPU) LDH_A_C() TCycles {
	offset := c.readR8(R8_C)
	value := c.bus.Read(0xFF00 + uint16(offset))
	c.writeR8(R8_A, value)
	return 8
}

func (c *CPU) LDH_CMEM_A() TCycles {
	offset := c.readR8(R8_C)
	value := c.readR8(R8_A)
	c.bus.Write(uint16(0xFF00)+uint16(offset), value)
	return 8
}

func (c *CPU) ADD_SP_E8() TCycles {
	offset := c.loadE8()
	oldSP := c.SP

	c.SP = uint16(int32(oldSP) + int32(offset))

	c.setZeroFlag(false)
	c.setSubtractionFlag(false)

	unsignedOffset := uint8(offset)
	c.setHalfCarryFlag(uint8(oldSP&0xF)+(unsignedOffset&0xF) > 0xF)
	c.setCarryFlag((oldSP&0xFF)+uint16(unsignedOffset) > 0xFF)
	return 16
}

func (c *CPU) JP_HL() TCycles {
	c.PC = c.readR16(R16_HL)
	return 4
}

func (c *CPU) LD_IMM16MEM_A() TCycles {
	address := c.loadImm16()
	value := c.readR8(R8_A)
	c.bus.Write(address, value)
	return 16
}

func (c *CPU) LD_A_IMM16MEM() TCycles {
	address := c.loadImm16()
	value := c.bus.Read(address)
	c.writeR8(R8_A, value)
	return 16
}

func (c *CPU) DI() TCycles {
	c.interruptsEnabled = false
	c.enableInterruptsNextOp = false
	return 4
}

func (c *CPU) EI() TCycles {
	c.enableInterruptsNextOp = true
	return 4
}

func (c *CPU) LD_HL_SP_PLUS_E8() TCycles {
	e8 := c.loadE8()
	oldSP := c.SP
	result := uint16(int32(oldSP) + int32(e8))
	c.writeR16(R16_HL, result)

	unsignedE8 := uint16(uint8(e8))
	c.setZeroFlag(false)
	c.setSubtractionFlag(false)
	c.setHalfCarryFlag((oldSP&0xF)+(unsignedE8&0xF) > 0xF)
	c.setCarryFlag((oldSP&0xFF)+unsignedE8 > 0xFF)
	return 12
}

func (c *CPU) LD_SP_HL() TCycles {
	value := c.readR16(R16_HL)
	c.writeR16(R16_SP, value)
	return 8
}

func (c *CPU) popStack() uint8 {
	value := c.bus.Read(c.SP)
	c.SP += 1
	return value
}

func (c *CPU) pushStack(value uint8) {
	c.SP -= 1
	c.bus.Write(c.SP, value)
}

func (c *CPU) pushStack16(value uint16) {
	high := uint8(value >> 8)
	low := uint8(value & 0xFF)

	c.pushStack(high)
	c.pushStack(low)
}

func (c *CPU) popStack16() uint16 {
	// Pop low byte first, then high byte
	low := uint16(c.popStack())
	high := uint16(c.popStack())
	return (high << 8) | low
}
