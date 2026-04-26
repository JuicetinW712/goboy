package cpu

type RegisterState struct {
	AF  uint16
	BC  uint16
	DE  uint16
	HL  uint16
	SP  uint16
	PC  uint16
	IME bool
}

func (c *CPU) GetRegisterState() RegisterState {
	return RegisterState{
		AF:  c.registers.getAF(),
		BC:  c.registers.getBC(),
		DE:  c.registers.getDE(),
		HL:  c.registers.getHL(),
		SP:  c.SP,
		PC:  c.PC,
		IME: c.interruptsEnabled,
	}
}

func (c *CPU) GetFlagState() string {
	// ZNHC
	baseStr := ""
	if c.getZeroFlag() {
		baseStr += "Z"
	} else {
		baseStr += "-"
	}
	if c.getSubtractionFlag() {
		baseStr += "N"
	} else {
		baseStr += "-"
	}
	if c.getHalfCarryFlag() {
		baseStr += "H"
	} else {
		baseStr += "-"
	}
	if c.getCarryFlag() {
		baseStr += "C"
	} else {
		baseStr += "-"
	}

	return baseStr
}
