package cpu

type R8 uint8
type R16 uint8
type R16Mem uint8
type R16Stk uint8
type COND uint8

const (
	R8_B R8 = iota
	R8_C
	R8_D
	R8_E
	R8_H
	R8_L
	R8_HL
	R8_A
	R8_F
)

const (
	R16_BC R16 = iota
	R16_DE
	R16_HL
	R16_SP
)

const (
	R16Mem_BC R16Mem = iota
	R16Mem_DE
	R16Mem_HLI
	R16Mem_HLD
)

const (
	R16Stk_BC R16Stk = iota
	R16Stk_DE
	R16Stk_HL
	R16Stk_AF
)

const (
	COND_NZ COND = iota
	COND_Z
	COND_NC
	COND_C
)

func (c *CPU) getR8FromIdx(registerNum uint8) R8 {
	switch registerNum {
	case 0:
		return R8_B
	case 1:
		return R8_C
	case 2:
		return R8_D
	case 3:
		return R8_E
	case 4:
		return R8_H
	case 5:
		return R8_L
	case 6:
		return R8_HL
	case 7:
		return R8_A
	default:
		panic("Invalid registerNum for getR8FromIdx")
	}
}

func (c *CPU) readR8(r R8) uint8 {
	switch r {
	case R8_B:
		return c.registers.B
	case R8_C:
		return c.registers.C
	case R8_D:
		return c.registers.D
	case R8_E:
		return c.registers.E
	case R8_H:
		return c.registers.H
	case R8_L:
		return c.registers.L
	case R8_HL:
		return c.bus.Read(c.registers.getHL())
	case R8_A:
		return c.registers.A
	case R8_F:
		return c.registers.F
	}
	panic("reading from invalid R8")
}

func (c *CPU) writeR8(r R8, val uint8) {
	switch r {
	case R8_B:
		c.registers.B = val
	case R8_C:
		c.registers.C = val
	case R8_D:
		c.registers.D = val
	case R8_E:
		c.registers.E = val
	case R8_H:
		c.registers.H = val
	case R8_L:
		c.registers.L = val
	case R8_HL:
		c.bus.Write(c.registers.getHL(), val)
	case R8_A:
		c.registers.A = val
	case R8_F:
		c.registers.F = val
	default:
		panic("writing to invalid R8")
	}
}

func (c *CPU) readR16(r R16) uint16 {
	switch r {
	case R16_BC:
		return c.registers.getBC()
	case R16_DE:
		return c.registers.getDE()
	case R16_HL:
		return c.registers.getHL()
	case R16_SP:
		return c.SP
	default:
		panic("reading from invalid R16")
	}
}

func (c *CPU) writeR16(r R16, val uint16) {
	switch r {
	case R16_BC:
		c.registers.setBC(val)
	case R16_DE:
		c.registers.setDE(val)
	case R16_HL:
		c.registers.setHL(val)
	case R16_SP:
		c.SP = val
	default:
		panic("writing to invalid R16")
	}
}

func (c *CPU) readR16Mem(r R16Mem) uint8 {
	addr := uint16(0xFEA0) // init to unused memory

	switch r {
	case R16Mem_BC:
		addr = c.readR16(R16_BC)
	case R16Mem_DE:
		addr = c.readR16(R16_DE)
	case R16Mem_HLI:
		addr = c.readR16(R16_HL)
		c.writeR16(R16_HL, addr+1)
	case R16Mem_HLD:
		addr = c.readR16(R16_HL)
		c.writeR16(R16_HL, addr-1)
	default:
		panic("reading from invalid R16Mem")
	}

	return c.bus.Read(addr)
}

func (c *CPU) writeR16Mem(r R16Mem, val uint8) {
	addr := uint16(0xFEA0) // init to unused memory

	switch r {
	case R16Mem_BC:
		addr = c.readR16(R16_BC)
	case R16Mem_DE:
		addr = c.readR16(R16_DE)
	case R16Mem_HLI:
		addr = c.readR16(R16_HL)
		c.writeR16(R16_HL, addr+1)
	case R16Mem_HLD:
		addr = c.readR16(R16_HL)
		c.writeR16(R16_HL, addr-1)
	default:
		panic("writing to invalid R16Mem")
	}

	c.bus.Write(addr, val)
}

func (c *CPU) readR16Stk(r R16Stk) uint16 {
	switch r {
	case R16Stk_BC:
		return c.readR16(R16_BC)
	case R16Stk_DE:
		return c.readR16(R16_DE)
	case R16Stk_HL:
		return c.readR16(R16_HL)
	case R16Stk_AF:
		return (uint16(c.registers.A) << 8) | (uint16(c.registers.F) & 0xF0)
	default:
		panic("reading invalid R16Stk")
	}
}

func (c *CPU) writeR16Stk(r R16Stk, value uint16) {
	switch r {
	case R16Stk_BC:
		c.writeR16(R16_BC, value)
	case R16Stk_DE:
		c.writeR16(R16_DE, value)
	case R16Stk_HL:
		c.writeR16(R16_HL, value)
	case R16Stk_AF:
		c.writeR8(R8_A, uint8(value>>8))
		c.writeR8(R8_F, uint8(value&0xF0)) // Need to mask lower 4 bits to ensure always 0
	default:
		panic("reading invalid R16Stk")
	}
}
