package cpu

type TCycles uint8

const (
	VBLANK_ADDR   uint16 = 0x0040
	LCD_STAT_ADDR uint16 = 0x0048
	TIMER_ADDR    uint16 = 0x0050
	SERIAL_ADDR   uint16 = 0x0058
	JOYPAD_ADDR   uint16 = 0x60
)

type Bus interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

type CPU struct {
	registers Registers
	PC        uint16 // Program Counter
	SP        uint16 // Stack Pointer

	bus                    Bus
	interruptsEnabled      bool // IME
	enableInterruptsNextOp bool // EI delay flag stage 1
	imeEnableRequested     bool // EI delay flag stage 2

	Halted  bool
	Stopped bool
	haltBug bool
}

func CreateCPU(bus Bus) *CPU {
	return &CPU{
		registers: CreateRegisters(),
		PC:        0x0100,
		SP:        0xFFFE,

		bus:                    bus,
		interruptsEnabled:      true,
		enableInterruptsNextOp: false,
		imeEnableRequested:     false,

		Halted:  false,
		Stopped: false,
	}
}

func (c *CPU) Step() TCycles {
	interruptFlags := c.bus.Read(0xFF0F)
	enabledInterrupts := c.bus.Read(0xFFFF)
	pendingInterrupts := interruptFlags & enabledInterrupts

	// Halted
	if c.Halted {
		if pendingInterrupts != 0 {
			c.Halted = false
		} else {
			return 4
		}
	}

	// Handle interrupts before instruction fetch if IME is set
	if c.interruptsEnabled && pendingInterrupts != 0 {
		return c.handleInterrupts(pendingInterrupts)
	}

	// Execute instruction
	opcode := c.bus.Read(c.PC)
	if !c.haltBug {
		c.PC++
	} else {
		c.haltBug = false
	}

	cycles := c.executeInstruction(opcode)

	// EI delay handling: IME is enabled AFTER the instruction following EI
	if c.imeEnableRequested {
		c.interruptsEnabled = true
		c.imeEnableRequested = false
	}
	if c.enableInterruptsNextOp {
		c.imeEnableRequested = true
		c.enableInterruptsNextOp = false
	}

	return cycles
}

func (c *CPU) executeInstruction(opcode uint8) TCycles {
	if opcode == 0xCB {
		cbOpcode := c.bus.Read(c.PC)
		c.PC += 1
		return c.executeBlockCB(cbOpcode)
	} else {
		block := opcode >> 6
		switch block {
		case 0:
			return c.executeBlock0(opcode)
		case 1:
			return c.executeBlock1(opcode)
		case 2:
			return c.executeBlock2(opcode)
		case 3:
			return c.executeBlock3(opcode)
		default:
			panic("Opcode reached invalid block")
		}
	}
}

func (c *CPU) handleInterrupts(pendingInterrupts uint8) TCycles {
	c.interruptsEnabled = false
	interruptFlags := c.bus.Read(0xFF0F)

	// Call in order of priority and clear flag
	// VBlank (0x40), LCD STAT (0x48), Timer (0x50), Serial (0x58), Joypad (0x60)

	var addr uint16
	var bit uint8

	if pendingInterrupts&0x01 != 0 {
		bit, addr = 0, VBLANK_ADDR
	} else if pendingInterrupts&0x02 != 0 {
		bit, addr = 1, LCD_STAT_ADDR
	} else if pendingInterrupts&0x04 != 0 {
		bit, addr = 2, TIMER_ADDR
	} else if pendingInterrupts&0x08 != 0 {
		bit, addr = 3, SERIAL_ADDR
	} else if pendingInterrupts&0x10 != 0 {
		bit, addr = 4, JOYPAD_ADDR
	}

	// Clear the flag in IF and call the address
	c.bus.Write(0xFF0F, interruptFlags & ^(1<<bit))
	c.call(addr)
	return 20
}

func (c *CPU) loadE8() int8 {
	val := int8(c.bus.Read(c.PC))
	c.PC += 1
	return val
}

func (c *CPU) loadImm8() uint8 {
	val := c.bus.Read(c.PC)
	c.PC += 1
	return val
}

func (c *CPU) loadImm16() uint16 {
	low := c.bus.Read(c.PC)
	c.PC += 1
	high := c.bus.Read(c.PC)
	c.PC += 1
	return uint16(high)<<8 | uint16(low)
}

func (c *CPU) getZeroFlag() bool {
	flags := c.readR8(R8_F)
	return (flags>>7)&1 == 1
}

func (c *CPU) getSubtractionFlag() bool {
	flags := c.readR8(R8_F)
	return (flags>>6)&1 == 1
}

func (c *CPU) getHalfCarryFlag() bool {
	flags := c.readR8(R8_F)
	return (flags>>5)&1 == 1
}

func (c *CPU) getCarryFlag() bool {
	flags := c.readR8(R8_F)
	return (flags>>4)&1 == 1
}

func (c *CPU) setZeroFlag(val bool) {
	flags := c.readR8(R8_F)
	if val {
		c.writeR8(R8_F, flags|0b10000000)
	} else {
		c.writeR8(R8_F, flags&0b01111111)
	}
}

func (c *CPU) setSubtractionFlag(val bool) {
	flags := c.readR8(R8_F)
	if val {
		c.writeR8(R8_F, flags|0b01000000)
	} else {
		c.writeR8(R8_F, flags&0b10111111)
	}
}

func (c *CPU) setHalfCarryFlag(val bool) {
	flags := c.readR8(R8_F)
	if val {
		c.writeR8(R8_F, flags|0b00100000)
	} else {
		c.writeR8(R8_F, flags&0b11011111)
	}
}

func (c *CPU) setCarryFlag(val bool) {
	flags := c.readR8(R8_F)
	if val {
		c.writeR8(R8_F, flags|0b00010000)
	} else {
		c.writeR8(R8_F, flags&0b11101111)
	}
}

func (c *CPU) evaluateCondition(condition COND) bool {
	switch condition {
	case COND_NZ:
		return !c.getZeroFlag()
	case COND_Z:
		return c.getZeroFlag()
	case COND_NC:
		return !c.getCarryFlag()
	case COND_C:
		return c.getCarryFlag()
	default:
		panic("Invalid COND")
	}
}
