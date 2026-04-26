package mbc1

const (
	// Bank sizes
	ROM_BANK_SIZE uint16 = 0x4000
	RAM_BANK_SIZE uint16 = 0x2000

	// Addresses that modify the MBC state
	RAM_ENABLE_START  uint16 = 0x0000
	RAM_ENABLE_END    uint16 = 0x1FFF
	ROM_LOWER_START   uint16 = 0x2000
	ROM_LOWER_END     uint16 = 0x3FFF
	ROM_UPPER_START   uint16 = 0x4000
	ROM_UPPER_END     uint16 = 0x5FFF
	MODE_SELECT_START uint16 = 0x6000
	MODE_SELECT_END   uint16 = 0x7FFF

	// ROM and RAM addresses
	ROM_BANK_0_START uint16 = 0x0000
	ROM_BANK_0_END   uint16 = 0x3FFF
	ROM_BANK_N_START uint16 = 0x4000
	ROM_BANK_N_END   uint16 = 0x7FFF
	RAM_BANK_START   uint16 = 0xA000
	RAM_BANK_END     uint16 = 0xBFFF
)

// 0x0000-0x1FFF: RAM Enable register
// 	4 bit register: 0x0A enables RAM, otherwise disable
// 0x2000-0x3FFF: ROM Bank register
//	5 bit register: Selects ROM Bank
// 0x4000-0x5FFF: RAM Bank register
// 	2 bit register: Selects RAM bank in 4 MBit/256 kbit and upper two ROM address in 16MBit/64 kbit
// 0x6000-0x7FFF: Mode register
//  1 bit register: 16 Mbit ROM/64 kbit SRAM mode ('0') and 4 Mbit ROM/256 kbit SRAM mode ('1')

type MBC_1 struct {
	rom []uint8
	ram []uint8

	ramEnabled      bool
	ramBankingMode  bool
	romBankRegister uint16
	ramBankRegister uint16

	numRomBanks uint16
	numRamBanks uint16
}

func CreateMBC1(rom []uint8, ram []uint8) *MBC_1 {
	nr := len(rom) / int(ROM_BANK_SIZE)
	na := len(ram) / int(RAM_BANK_SIZE)

	// Ensure at least 1 ROM Bank to prevent division by zero
	if nr == 0 {
		nr = 1
	}

	return &MBC_1{
		rom:             rom,
		ram:             ram,
		romBankRegister: 1,
		ramBankRegister: 0,
		numRomBanks:     uint16(nr),
		numRamBanks:     uint16(na),
	}
}

func (mbc *MBC_1) Read(address uint16) uint8 {
	switch {
	case address <= ROM_BANK_0_END:
		romBank := uint16(0)

		// When you write to the romBankRegister, 0 becomes 1
		// so 0x20, 0x40, ..., not writable.
		// To address this, when in ramBankingMode, we
		if mbc.ramBankingMode {
			romBank = mbc.ramBankRegister << 5
		}

		romBank %= mbc.numRomBanks
		return mbc.rom[uint32(address)+(uint32(romBank)*uint32(ROM_BANK_SIZE))]
	case address <= ROM_BANK_N_END:
		var romBank uint16

		// If romBanking mode, use ramBankRegister as upper 2 bits
		if !mbc.ramBankingMode {
			romBank = (mbc.ramBankRegister << 5) | mbc.romBankRegister
		} else {
			romBank = mbc.romBankRegister
		}

		// Can't map to 0 on romBankRegister
		if mbc.romBankRegister == 0 {
			romBank += 1
		}

		romBank %= mbc.numRomBanks
		offset := uint16(address - ROM_BANK_N_START)
		return mbc.rom[uint32(offset)+(uint32(romBank)*uint32(ROM_BANK_SIZE))]
	case address <= RAM_BANK_END:
		// No offset if ram banking mode is not turned on
		if !mbc.ramEnabled || mbc.numRamBanks == 0 {
			return 0xFF
		}

		ramBank := uint16(0)
		if mbc.ramBankingMode {
			ramBank = mbc.ramBankRegister
		}

		ramBank %= mbc.numRamBanks
		offset := uint16(address - RAM_BANK_START)
		return mbc.ram[uint32(offset)+(uint32(ramBank)*uint32(RAM_BANK_SIZE))]
	default:
		return 0xFF
	}
}

func (mbc *MBC_1) Write(address uint16, value uint8) {
	switch {
	case address <= RAM_ENABLE_END:
		// Writing 0xA enables RAM while writing anything else disables it
		mbc.ramEnabled = ((value & 0xF) == 0xA)

	case address <= ROM_LOWER_END:
		mbc.romBankRegister = uint16(value & 0x1F)

	case address <= ROM_UPPER_END:
		mbc.ramBankRegister = uint16(value & 0x03)

	case address <= MODE_SELECT_END:
		mbc.ramBankingMode = (value & 0b1) == 1

	case address <= RAM_BANK_END:
		if mbc.ramEnabled && mbc.numRamBanks > 0 {
			ramBank := uint16(0)

			if mbc.ramBankingMode {
				ramBank = mbc.ramBankRegister
			}

			ramBank %= mbc.numRamBanks
			offset := uint16(address - RAM_BANK_START)
			mbc.ram[uint32(offset)+(uint32(ramBank)*uint32(RAM_BANK_SIZE))] = value
		}
	default:
		panic("MBC should not be handling address for writes")
	}
}
