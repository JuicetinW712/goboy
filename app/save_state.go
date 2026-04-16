package app

import (
	"encoding/gob"
	"fmt"
	"os"

	"goboy/app/bus"
	"goboy/app/cartridge"
	"goboy/app/cartridge/mbc1"
	"goboy/app/cartridge/mbc3"
	"goboy/app/cartridge/mbc5"
	"goboy/app/cartridge/rom_only"
	"goboy/app/cpu"
	"goboy/app/io"
	"goboy/app/ppu"
	"goboy/app/timer"
)

func init() {
	gob.Register(mbc1.MBC1State{})
	gob.Register(mbc3.MBC3State{})
	gob.Register(mbc5.MBC5State{})
	gob.Register(rom_only.ROMOnlyState{})
}

type SaveState struct {
	BusState   bus.MMUState
	CPUState   cpu.CPUState
	CartState  cartridge.CartridgeState
	IOState    io.IOState
	PPUState   ppu.PPUState
	TimerState timer.TimerState
}

func (g *Game) SaveState() error {
	file, err := os.Create("game.save")
	if err != nil {
		return fmt.Errorf("Failed to create file")
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)

	saveState := SaveState{
		BusState:   g.bus.GetState(),
		CPUState:   g.cpu.GetState(),
		CartState:  g.cartridge.GetState(),
		IOState:    g.io.GetState(),
		PPUState:   g.ppu.GetState(),
		TimerState: g.timer.GetState(),
	}

	err = encoder.Encode(saveState)
	if err != nil {
		return fmt.Errorf("Failed to encode emulator state: %w", err)
	}

	return nil
}

func (g *Game) LoadSaveState(saveState SaveState) error {
	if err := g.cartridge.LoadState(saveState.CartState); err != nil {
		return fmt.Errorf("Failed to load cartridge state: %w", err)
	}
	g.ppu.LoadState(&saveState.PPUState)
	g.timer.LoadState(saveState.TimerState)
	g.io.LoadState(saveState.IOState)
	g.bus.LoadState(&saveState.BusState)
	g.cpu.LoadState(&saveState.CPUState)
	return nil
}

func LoadSaveFromFile(fileName string) (SaveState, error) {
	var saveState SaveState

	file, err := os.Open(fileName)
	if err != nil {
		return SaveState{}, fmt.Errorf("Failed to read save file")
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	if err := decoder.Decode(&saveState); err != nil {
		return SaveState{}, fmt.Errorf("Failed to decode save file into struct: %w", err)
	}

	return saveState, nil
}
