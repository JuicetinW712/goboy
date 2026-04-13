package mbc3

const (
	// Bank sizes
	ROM_BANK_SIZE uint16 = 0x4000
	RAM_BANK_SIZE uint16 = 0x2000

	// Registers
	RAM_AND_TIMER_ENABLE_START uint16 = 0x0000
	RAM_AND_TIMER_ENABLE_END   uint16 = 0x1FFF
	ROM_LOWER_START            uint16 = 0x2000
	ROM_LOWER_END              uint16 = 0x3FFF
	ROM_UPPER_START            uint16 = 0x4000
	ROM_UPPER_END              uint16 = 0x5FFF
	MODE_SELECT_START          uint16 = 0x6000
	MODE_SELECT_END            uint16 = 0x7FFF

	// ROM and RAM addresses
	ROM_BANK_0_START uint16 = 0x0000
	ROM_BANK_0_END   uint16 = 0x3FFF
	ROM_BANK_N_START uint16 = 0x4000
	ROM_BANK_N_END   uint16 = 0x7FFF
	RAM_BANK_START   uint16 = 0xA000
	RAM_BANK_END     uint16 = 0xBFFF
)

// TODO - implement timer instead of using returning 0xFF

// In MBC1, there is logic to increment the lower bits of the ROM Bank Number
// if set to 0 since bank 0 is another address. This led to a problem where
// 0x20, 0x40, 0x60, .... is not reachable since only the upper bits are written to
// and the lower bits are all 0 so the hardware increments it

// With MBC3, the upper and lower bits are together so this is not an issue

// 0x0000-0x1FFF: RAM Enable register
// 	4 bit register: 0x0A enables RAM, otherwise disable
// 0x2000-0x3FFF: ROM Bank register
//	7 bit register: Selects ROM Bank
// 0x4000-0x5FFF: RAM Bank register
//  00-07: Selects rom bank
//  08-0C  Corresponding RTC Register
// 		08: RTC_S (0-59), seconds
// 		09: RTC_M (0-59), minutes
// 		0A: RTC_H (0-23), hours
// 		0B: RTC_DL (0-FF), lower 8 bits of day
// 		0C: RTC_DH, ()

// 	2 bit register: Selects RAM bank in 4 MBit/256 kbit and upper two ROM address in 16MBit/64 kbit
// 0x6000-0x7FFF: Mode register
//  1 bit register: 16 Mbit ROM/64 kbit SRAM mode ('0') and 4 Mbit ROM/256 kbit SRAM mode ('1')

type MBC_3 struct {
	rom []uint8
	ram []uint8

	ramAndTimerEnabled bool
	romBankRegister    uint16
	ramAndRtcRegister  uint16

	numRomBanks uint16
	numRamBanks uint16
}

func CreateMBC3(rom []uint8, ram []uint8) *MBC_3 {
	nr := len(rom) / int(ROM_BANK_SIZE)
	na := len(ram) / int(RAM_BANK_SIZE)

	// Ensure at least 1 ROM Bank to prevent division by zero
	if nr == 0 {
		nr = 1
	}

	return &MBC_3{
		rom:               rom,
		ram:               ram,
		romBankRegister:   1,
		ramAndRtcRegister: 0,
		numRomBanks:       uint16(nr),
		numRamBanks:       uint16(na),
	}
}

func (mbc *MBC_3) Read(address uint16) uint8 {
	switch {
	case address <= ROM_BANK_0_END:
		return mbc.rom[address]
	case address <= ROM_BANK_N_END:
		romBank := mbc.romBankRegister

		// Can't map to 0 on romBankRegister
		if mbc.romBankRegister == 0 {
			romBank += 1
		}

		romBank %= mbc.numRomBanks
		offset := uint16(address - ROM_BANK_N_START)
		return mbc.rom[uint32(offset)+(uint32(romBank)*uint32(ROM_BANK_SIZE))]
	case address <= RAM_BANK_END:
		// No offset if ram banking mode is not turned on
		if !mbc.ramAndTimerEnabled || mbc.numRamBanks == 0 {
			return 0xFF
		}

		// Below 7, rom bank, above 7 rtc register
		if mbc.ramAndRtcRegister <= 7 {
			ramBank := mbc.ramAndRtcRegister
			ramBank %= mbc.numRamBanks
			offset := uint16(address - RAM_BANK_START)
			return mbc.ram[uint32(offset)+(uint32(ramBank)*uint32(RAM_BANK_SIZE))]
		} else {
			// 0xFF for now
			return 0xFF
		}
	default:
		return 0xFF
	}
}

func (mbc *MBC_3) Write(address uint16, value uint8) {
	switch {
	case address <= RAM_AND_TIMER_ENABLE_END:
		mbc.ramAndTimerEnabled = ((value & 0xF) == 0xA)
	case address <= 0x3FFF:
		mbc.romBankRegister = uint16(value & 0x7F)
	case address <= 0x5FFF:
		mbc.ramAndRtcRegister = uint16(value & 0xF)
	case address <= 0x7FFF:
		// Latch clock data

		// Writing 0 and then 1, stores current time in RTC registers
	case address >= 0xA000 && address <= 0xBFFF:
		if !mbc.ramAndTimerEnabled {
			return
		}

		// Below 7, rom bank, above 7 rtc register
		if mbc.ramAndRtcRegister <= 7 {
			ramBank := mbc.ramAndRtcRegister
			ramBank %= mbc.numRamBanks
			offset := uint16(address - RAM_BANK_START)
			mbc.ram[uint32(offset)+(uint32(ramBank)*uint32(RAM_BANK_SIZE))] = value
		} else {

		}
	}
}
