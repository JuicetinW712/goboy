package bus

type MemoryRegion uint8

const (
	REGION_ROM_0 MemoryRegion = iota
	REGION_ROM_N
	REGION_VRAM
	REGION_ERAM
	REGION_WRAM
	REGION_ECHO_RAM
	REGION_OAM
	REGION_IO
	REGION_HRAM
	REGION_INTERRUPT_REGISTER
	REGION_UNUSABLE // address that shouldn't be used
	REGION_UNKNOWN  // impossible memory address
)

const (
	ROM_0_START        uint16 = 0x0000
	ROM_0_END          uint16 = 0x3FFF
	ROM_N_START        uint16 = 0x4000
	ROM_N_END          uint16 = 0x7FFF
	VRAM_START         uint16 = 0x8000
	VRAM_END           uint16 = 0x9FFF
	ERAM_START         uint16 = 0xA000
	ERAM_END           uint16 = 0xBFFF
	WRAM_START         uint16 = 0xC000
	WRAM_END           uint16 = 0xDFFF
	ECHO_RAM_START     uint16 = 0xE000
	ECHO_RAM_END       uint16 = 0xFDFF
	OAM_START          uint16 = 0xFE00
	OAM_END            uint16 = 0xFE9F
	UNUSABLE_START     uint16 = 0xFEA0
	UNUSABLE_END       uint16 = 0xFEFF
	IO_START           uint16 = 0xFF00
	IO_END             uint16 = 0xFF7F
	HRAM_START         uint16 = 0xFF80
	HRAM_END           uint16 = 0xFFFE
	INTERRUPT_REGISTER uint16 = 0xFFFF
)
