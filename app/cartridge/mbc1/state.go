package mbc1

import "errors"

type MBC1State struct {
	RamEnabled      bool
	RamBankingMode  bool
	RomBankRegister uint16
	RamBankRegister uint16
}

func (mbc *MBC_1) GetState() any {
	return MBC1State{
		RamEnabled:      mbc.ramEnabled,
		RamBankingMode:  mbc.ramBankingMode,
		RomBankRegister: mbc.romBankRegister,
		RamBankRegister: mbc.ramBankRegister,
	}
}

func (mbc *MBC_1) LoadState(state any) error {
	mbcState, ok := state.(MBC1State)
	if !ok {
		return errors.New("Failed to load MBC1 state")
	}
	mbc.ramEnabled = mbcState.RamEnabled
	mbc.ramBankingMode = mbcState.RamBankingMode
	mbc.romBankRegister = mbcState.RomBankRegister
	mbc.ramBankRegister = mbcState.RamBankRegister
	return nil
}
