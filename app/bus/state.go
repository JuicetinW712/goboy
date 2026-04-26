package bus

type MMUState struct {
	Wram              [8192]uint8
	Hram              [127]uint8
	InterruptRegister uint8
}

func (mmu *MMU) GetState() MMUState {
	return MMUState{
		Wram:              mmu.wram,
		Hram:              mmu.hram,
		InterruptRegister: mmu.interruptRegister,
	}
}

func (mmu *MMU) LoadState(state *MMUState) {
	mmu.wram = state.Wram
	mmu.hram = state.Hram
	mmu.interruptRegister = state.InterruptRegister
}
