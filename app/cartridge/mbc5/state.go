package mbc5

import "errors"

type MBC5State struct {
	RamEnabled      bool
	RomBankRegister uint16
	RamBankRegister uint8
}

func (mbc *MBC_5) GetState() any {
	return MBC5State{
		RamEnabled:      mbc.ramEnabled,
		RomBankRegister: mbc.romBankRegister,
		RamBankRegister: mbc.ramBankRegister,
	}
}

func (mbc *MBC_5) LoadState(state any) error {
	mbcState, ok := state.(MBC5State)
	if !ok {
		return errors.New("Failed to load MBC5 state")
	}
	mbc.ramEnabled = mbcState.RamEnabled
	mbc.romBankRegister = mbcState.RomBankRegister
	mbc.ramBankRegister = mbcState.RamBankRegister
	return nil
}
