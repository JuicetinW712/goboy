package mbc3

import "errors"

type MBC3State struct {
	RamAndTimerEnabled bool
	RomBankRegister    uint16
	RamAndRtcRegister  uint16
}

func (mbc *MBC_3) GetState() any {
	return MBC3State{
		RamAndTimerEnabled: mbc.ramAndTimerEnabled,
		RomBankRegister:    mbc.romBankRegister,
		RamAndRtcRegister:  mbc.ramAndRtcRegister,
	}
}

func (mbc *MBC_3) LoadState(state any) error {
	mbcState, ok := state.(MBC3State)
	if !ok {
		return errors.New("Failed to load MBC3 state")
	}
	mbc.ramAndTimerEnabled = mbcState.RamAndTimerEnabled
	mbc.romBankRegister = mbcState.RomBankRegister
	mbc.ramAndRtcRegister = mbcState.RamAndRtcRegister
	return nil
}
