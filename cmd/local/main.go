package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"goboy/app/bus"
	"goboy/app/cartridge"
	"goboy/app/cpu"
	"goboy/app/io"
	"goboy/app/joypad"
	"goboy/app/ppu"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	cpu    *cpu.CPU
	bus    *bus.MMU
	timer  *io.Timer
	joypad *joypad.Joypad
	ppu    *ppu.PPU
	io     *io.IO

	debugOn bool
	sync    bool
}

func (g *Game) Update() error {
	const cyclesPerFrame = 70224
	g.updateInput()

	cyclesRan := uint32(0)
	for cyclesRan < cyclesPerFrame {
		cycles := g.cpu.Step()

		for i := 0; i < int(cycles); i++ {
			g.timer.Tick()
			g.ppu.Step(1)
		}

		cyclesRan += uint32(cycles)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Wait until entire frame is drawn if sync is enabled
	// Prevents screen tearing
	if g.sync {
		screen.WritePixels(g.ppu.Image.Pix)
	} else {
		screen.WritePixels(g.ppu.Buffer.Pix)
	}

	if g.debugOn {
		debugInfo := fmt.Sprintf("FPS: %0.2f\nTPS: %0.2f\n",
			ebiten.ActualFPS(),
			ebiten.ActualTPS(),
		)
		ebitenutil.DebugPrint(screen, debugInfo)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 160, 144
}

func (g *Game) updateInput() {
	g.updateKeyState(ebiten.KeyW, joypad.DPAD_UP)
	g.updateKeyState(ebiten.KeyA, joypad.DPAD_LEFT)
	g.updateKeyState(ebiten.KeyS, joypad.DPAD_DOWN)
	g.updateKeyState(ebiten.KeyD, joypad.DPAD_RIGHT)
	g.updateKeyState(ebiten.KeyZ, joypad.BUTTON_A)
	g.updateKeyState(ebiten.KeyX, joypad.BUTTON_B)
	g.updateKeyState(ebiten.KeyEnter, joypad.BUTTON_START)
	g.updateKeyState(ebiten.KeyShift, joypad.BUTTON_SELECT)
}

func (g *Game) updateKeyState(key ebiten.Key, mapping joypad.Button) {
	if ebiten.IsKeyPressed(key) {
		g.joypad.PressButton(mapping)
	} else {
		g.joypad.ReleaseButton(mapping)
	}
}

func main() {
	// Parse flags
	scale := flag.Int("scale", 5, "Multiplier for screen size")
	debug := flag.Bool("debug", false, "Prints debug info to screen")
	noSync := flag.Bool("nosync", false, "Will draw to screen even if frame not fully rendered (possible screen tearing)")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: goboy [flags] <rom_file>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Create everything
	cart, err := cartridge.NewCartridgeFromFile(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	p := ppu.CreatePPU()
	t := &io.Timer{}
	j := joypad.CreateJoypad()
	i := io.NewIO(t, p, j)
	b := bus.NewMMU(cart, p, i)
	c := cpu.CreateCPU(b)

	game := &Game{
		cpu:     c,
		bus:     b,
		timer:   t,
		joypad:  j,
		ppu:     p,
		io:      i,
		debugOn: *debug,
		sync:    !*noSync,
	}

	// Connect interrupts
	p.OnVBlank = func() {
		i.RequestInterrupt(0)
	}
	p.OnStat = func() {
		i.RequestInterrupt(1)
	}
	t.RequestInterrupt = func(bit uint8) {
		i.RequestInterrupt(bit)
	}
	p.OnDMA = func(value uint8) {
		sourceBase := uint16(value) << 8

		for offset := range uint16(0xA0) {
			data := b.Read(sourceBase + offset)
			b.Write(0xFE00+offset, data)
		}
	}

	print(cart.String())
	ebiten.SetWindowSize(160**scale, 144**scale)
	ebiten.SetWindowTitle("GoBoy")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
