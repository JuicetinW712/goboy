package cartridge

// MBC is the memory bank controller that controls external
// RAM and ROM on a cartridge

type MBC interface {
	Read(address uint16) uint8
	Write(address uint16, value uint8)
}
