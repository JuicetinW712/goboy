package cpu

type CPUState struct {
	RegisterState RegisterState

	EnableInterruptsNextOp bool
	ImeEnableRequested     bool

	Halted  bool
	Stopped bool
	HaltBug bool
}

func (c *CPU) GetState() CPUState {
	return CPUState{
		RegisterState: c.GetRegisterState(),

		EnableInterruptsNextOp: c.enableInterruptsNextOp,
		ImeEnableRequested:     c.imeEnableRequested,

		Halted:  c.Halted,
		Stopped: c.Stopped,
		HaltBug: c.haltBug,
	}
}

func (c *CPU) LoadState(state *CPUState) {
	c.registers = Registers{
		A: uint8(state.RegisterState.AF >> 8),
		F: uint8(state.RegisterState.AF & 0xFF),
		B: uint8(state.RegisterState.BC >> 8),
		C: uint8(state.RegisterState.BC & 0xFF),
		D: uint8(state.RegisterState.DE >> 8),
		E: uint8(state.RegisterState.DE & 0xFF),
		H: uint8(state.RegisterState.HL >> 8),
		L: uint8(state.RegisterState.HL & 0xFF),
	}
	c.PC = state.RegisterState.PC
	c.SP = state.RegisterState.SP
	c.interruptsEnabled = state.RegisterState.IME
	c.enableInterruptsNextOp = state.EnableInterruptsNextOp
	c.imeEnableRequested = state.ImeEnableRequested
	c.Halted = state.Halted
	c.Stopped = state.Stopped
	c.haltBug = state.HaltBug
}
