package mbc5

const (
	ROM_BANK_SIZE uint32 = 0x4000
	RAM_BANK_SIZE uint32 = 0x2000

	RAM_ENABLE_END    uint16 = 0x1FFF
	ROM_BANK_LOW_END  uint16 = 0x2FFF // 0x2000-0x2FFF: Low 8 bits of ROM Bank
	ROM_BANK_HIGH_END uint16 = 0x3FFF // 0x3000-0x3FFF: 9th bit of ROM Bank
	RAM_BANK_SEL_END  uint16 = 0x5FFF // 0x4000-0x5FFF: RAM Bank (4 bits)

	ROM_BANK_0_END   uint16 = 0x3FFF
	ROM_BANK_N_START uint16 = 0x4000
	ROM_BANK_N_END   uint16 = 0x7FFF
	RAM_BANK_START   uint16 = 0xA000
	RAM_BANK_END     uint16 = 0xBFFF
)

type MBC_5 struct {
	rom []uint8
	ram []uint8

	ramEnabled      bool
	romBankRegister uint16 // Needs 9 bits
	ramBankRegister uint8  // Needs 4 bits

	numRomBanks uint16
	numRamBanks uint16
}

func CreateMBC5(rom []uint8, ram []uint8) *MBC_5 {
	nr := uint16(len(rom) / int(ROM_BANK_SIZE))
	na := uint16(len(ram) / int(RAM_BANK_SIZE))

	// Ensure at least 1 bank exists to avoid modulo by zero
	if nr == 0 {
		nr = 1
	}

	return &MBC_5{
		rom:             rom,
		ram:             ram,
		romBankRegister: 1,
		ramBankRegister: 0,
		numRomBanks:     nr,
		numRamBanks:     na,
	}
}

func (mbc *MBC_5) Read(address uint16) uint8 {
	switch {
	case address <= ROM_BANK_0_END:
		// MBC5 Fixed Range: Always Bank 0
		return mbc.rom[address]

	case address <= ROM_BANK_N_END:
		bank := uint32(mbc.romBankRegister % mbc.numRomBanks)
		offset := uint32(address - ROM_BANK_N_START)
		return mbc.rom[(bank*ROM_BANK_SIZE)+offset]

	case address >= RAM_BANK_START && address <= RAM_BANK_END:
		if !mbc.ramEnabled || mbc.numRamBanks == 0 {
			return 0xFF
		}
		bank := uint32(uint16(mbc.ramBankRegister) % mbc.numRamBanks)
		offset := uint32(address - RAM_BANK_START)
		return mbc.ram[(bank*RAM_BANK_SIZE)+offset]

	default:
		return 0xFF
	}
}

func (mbc *MBC_5) Write(address uint16, value uint8) {
	switch {
	case address <= RAM_ENABLE_END:
		mbc.ramEnabled = (value & 0x0F) == 0x0A

	case address <= ROM_BANK_LOW_END:
		// Set lower 8 bits, keep the 9th bit
		mbc.romBankRegister = (mbc.romBankRegister & 0x0100) | uint16(value)

	case address <= ROM_BANK_HIGH_END:
		// Set the 9th bit (bit 0 of the value)
		if value&0x01 != 0 {
			mbc.romBankRegister |= 0x0100
		} else {
			mbc.romBankRegister &= 0x00FF
		}

	case address <= RAM_BANK_SEL_END:
		// MBC5 supports up to 16 RAM banks (4 bits)
		mbc.ramBankRegister = value & 0x0F

	case address >= RAM_BANK_START && address <= RAM_BANK_END:
		if mbc.ramEnabled && mbc.numRamBanks > 0 {
			bank := uint32(uint16(mbc.ramBankRegister) % mbc.numRamBanks)
			offset := uint32(address - RAM_BANK_START)
			mbc.ram[(bank*RAM_BANK_SIZE)+offset] = value
		}
	}
}
