package bus

type MemoryDevice interface {
	Read(addr uint16) uint8
	Write(addr uint16, val uint8)
}

type MMU struct {
	cartridge         MemoryDevice
	wram              [8192]uint8
	hram              [127]uint8
	interruptRegister uint8

	ppu MemoryDevice
	io  MemoryDevice
}

func NewMMU(cartridge MemoryDevice, ppu MemoryDevice, io MemoryDevice) *MMU {
	return &MMU{
		cartridge: cartridge,
		ppu:       ppu,
		io:        io,
	}
}

func (mmu *MMU) getMemoryRegion(addr uint16) MemoryRegion {
	switch {
	case addr <= ROM_0_END:
		return REGION_ROM_0
	case addr <= ROM_N_END:
		return REGION_ROM_N
	case addr <= VRAM_END:
		return REGION_VRAM
	case addr <= ERAM_END:
		return REGION_ERAM
	case addr <= WRAM_END:
		return REGION_WRAM
	case addr <= ECHO_RAM_END:
		return REGION_ECHO_RAM
	case addr <= OAM_END:
		return REGION_OAM
	case addr <= UNUSABLE_END:
		return REGION_UNUSABLE
	case addr <= IO_END:
		return REGION_IO
	case addr <= HRAM_END:
		return REGION_HRAM
	case addr == INTERRUPT_REGISTER:
		return REGION_INTERRUPT_REGISTER
	default:
		return REGION_UNKNOWN
	}
}

func (mmu *MMU) Read(addr uint16) uint8 {
	region := mmu.getMemoryRegion(addr)

	switch region {
	case REGION_ROM_0, REGION_ROM_N, REGION_ERAM:
		return mmu.cartridge.Read(addr)
	case REGION_VRAM:
		return mmu.ppu.Read(addr)
	case REGION_WRAM:
		return mmu.wram[addr-WRAM_START]
	case REGION_ECHO_RAM:
		return mmu.wram[addr-ECHO_RAM_START]
	case REGION_OAM:
		return mmu.ppu.Read(addr)
	case REGION_IO:
		return mmu.io.Read(addr)
	case REGION_HRAM:
		return mmu.hram[addr-HRAM_START]
	case REGION_INTERRUPT_REGISTER:
		return mmu.interruptRegister
	case REGION_UNUSABLE:
		return 0xFF
	default:
		return 0xFF
	}
}

func (mmu *MMU) Write(addr uint16, val uint8) {
	region := mmu.getMemoryRegion(addr)

	switch region {
	case REGION_ROM_0, REGION_ROM_N, REGION_ERAM:
		mmu.cartridge.Write(addr, val)
	case REGION_VRAM:
		mmu.ppu.Write(addr, val)
	case REGION_WRAM:
		mmu.wram[addr-WRAM_START] = val
	case REGION_ECHO_RAM:
		mmu.wram[addr-ECHO_RAM_START] = val
	case REGION_OAM:
		mmu.ppu.Write(addr, val)
	case REGION_IO:
		mmu.io.Write(addr, val)
	case REGION_HRAM:
		mmu.hram[addr-HRAM_START] = val
	case REGION_INTERRUPT_REGISTER:
		mmu.interruptRegister = val
	default:
	}
}
