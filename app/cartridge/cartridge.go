package cartridge

import (
	"bytes"
	"fmt"
	"os"
)

const (
	KILOBYTE = 1024

	TITLE_START_ADDR = 0x0134
	TITLE_END_ADDR   = 0x0143
	MBC_TYPE_ADDR    = 0x0147
	RAM_SIZE_ADDR    = 0x0149
)

type Cartridge struct {
	rom      []uint8
	ram      []uint8
	title    string
	filename string

	mbcType MBCType
	mbc     MBC
}

func NewCartridge(filename string, rom []uint8) *Cartridge {
	romCopy := make([]uint8, len(rom))
	copy(romCopy, rom)

	rawTitle := rom[TITLE_START_ADDR:TITLE_END_ADDR]
	title := string(bytes.Trim(rawTitle, "\x00"))

	ramSize := getRamSize(rom[RAM_SIZE_ADDR])
	ram := make([]uint8, ramSize)
	mbcType := getMBCType(rom[MBC_TYPE_ADDR])

	return &Cartridge{
		rom:      romCopy,
		ram:      ram,
		title:    title,
		filename: filename,

		mbcType: mbcType,
		mbc:     getMBC(mbcType, romCopy, ram),
	}
}

func NewCartridgeFromFile(filename string) (*Cartridge, error) {
	rom, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return NewCartridge(filename, rom), nil
}

func (c *Cartridge) Read(addr uint16) uint8 {
	return c.mbc.Read(addr)
}

func (c *Cartridge) Write(addr uint16, val uint8) {
	c.mbc.Write(addr, val)
}

func (c *Cartridge) String() string {
	return fmt.Sprintf(
		"---Cartridge Info---\n"+
			"Title:    %s\n"+
			"MBC Type: %d\n"+
			"ROM Size: %d KB\n"+
			"RAM Size: %d KB\n",
		c.title,
		c.mbcType,
		len(c.rom)/KILOBYTE,
		len(c.ram)/KILOBYTE,
	)
}
