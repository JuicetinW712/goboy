package cartridge

import (
	"goboy/app/cartridge/mbc1"
	"goboy/app/cartridge/mbc3"
	"goboy/app/cartridge/mbc5"
	"goboy/app/cartridge/rom_only"
)

type MBCType uint8

const (
	TYPE_ROM_ONLY MBCType = iota
	TYPE_MBC1
	TYPE_MBC2
	TYPE_MBC3
	TYPE_MBC4
	TYPE_MBC5
	TYPE_UNSUPPORTED
)

func getMBCType(code uint8) MBCType {
	switch code {
	case 0x00:
		return TYPE_ROM_ONLY
	case 0x01:
		return TYPE_MBC1
	case 0x0F, 0x10, 0x11, 0x12, 0x13:
		return TYPE_MBC3
	case 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E:
		return TYPE_MBC5
	default:
		return TYPE_UNSUPPORTED
	}
}

func getMBC(mbcType MBCType, rom []uint8, ram []uint8) MBC {
	switch mbcType {
	case TYPE_ROM_ONLY:
		return rom_only.CreateROMOnlyMBC(rom)
	case TYPE_MBC1:
		return mbc1.CreateMBC1(rom, ram)
	case TYPE_MBC3:
		return mbc3.CreateMBC3(rom, ram)
	case TYPE_MBC5:
		return mbc5.CreateMBC5(rom, ram)
	default:
		panic("Unsupported MBC Type")
	}
}

func getRamSize(code uint8) uint {
	switch code {
	case 0x00:
		return 0
	case 0x02:
		return 8 * KILOBYTE
	case 0x03:
		return 32 * KILOBYTE
	case 0x04:
		return 128 * KILOBYTE
	case 0x05:
		return 64 * KILOBYTE
	default:
		panic("Invalid RAM Size Code")
	}
}
