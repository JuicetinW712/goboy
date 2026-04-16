package app

import (
	"fmt"
	"image/color"

	"goboy/app/bus"
	"goboy/app/cartridge"
	"goboy/app/cpu"
	"goboy/app/io"
	"goboy/app/joypad"
	"goboy/app/ppu"
	"goboy/app/timer"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	SCREEN_WIDTH      = 160
	SCREEN_HEIGHT     = 144
	DEBUG_PANEL_WIDTH = 144
)

type Game struct {
	cpu    *cpu.CPU
	bus    *bus.MMU
	timer  *timer.Timer
	joypad *joypad.Joypad
	ppu    *ppu.PPU
	io     *io.IO

	debugOn bool
	sync    bool
	pause   bool
	speed   uint32

	originalFilterEnabled bool
}

func NewGame(cart *cartridge.Cartridge, debug bool, sync bool, scale int) *Game {
	p := ppu.CreatePPU()
	t := &timer.Timer{}
	j := joypad.CreateJoypad()
	i := io.NewIO(t, p, j)
	b := bus.NewMMU(cart, p, i)
	c := cpu.CreateCPU(b)

	game := &Game{
		cpu:                   c,
		bus:                   b,
		timer:                 t,
		joypad:                j,
		ppu:                   p,
		io:                    i,
		debugOn:               debug,
		sync:                  sync,
		speed:                 1,
		originalFilterEnabled: false,
	}

	// Pass in functions to handle interrupts
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

	if game.debugOn {
		ebiten.SetWindowSize((SCREEN_WIDTH+DEBUG_PANEL_WIDTH)*scale, (SCREEN_HEIGHT)*scale)
	} else {
		ebiten.SetWindowSize((SCREEN_WIDTH)*scale, (SCREEN_HEIGHT)*scale)
	}

	return game
}

func (g *Game) Update() error {
	const cyclesPerFrame = 70224
	g.updateInput()
	g.updatePauseState()
	g.updateSpeed()
	g.updateFilter()

	if g.pause {
		return nil
	}

	cyclesRan := uint32(0)
	for cyclesRan < (cyclesPerFrame * g.speed) {
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
	g.drawGame(screen)

	if g.debugOn {
		g.drawDebugPanel(screen)
	}
}

func (g *Game) drawGame(screen *ebiten.Image) {
	if g.debugOn {
		game := ebiten.NewImage(SCREEN_WIDTH, SCREEN_HEIGHT)

		if g.sync {
			game.WritePixels(g.ppu.Image.Pix)
		} else {
			game.WritePixels(g.ppu.Buffer.Pix)
		}

		opt := &ebiten.DrawImageOptions{}
		if g.originalFilterEnabled {
			opt.ColorScale.Scale(0.608, 0.737, 0.059, 1.0)
		}

		screen.DrawImage(game, opt)
	} else {
		if g.sync {
			screen.WritePixels(g.ppu.Image.Pix)
		} else {
			screen.WritePixels(g.ppu.Buffer.Pix)
		}
	}
}

func (g *Game) drawDebugPanel(screen *ebiten.Image) {
	debugPanel := ebiten.NewImage(DEBUG_PANEL_WIDTH, SCREEN_HEIGHT)

	registerState := g.cpu.GetRegisterState()
	flagState := g.cpu.GetFlagState()
	timerState := g.timer.GetTimerState()

	stats := fmt.Sprintf(
		"FPS: %0.2f TPS: %0.2f\nSpeed: %d\nAF:  %04X PC:  %04X\nBC:  %04X SP:  %04X\nDE:  %04X HL:  %04X\nIME: %t\n"+
			"FLAGS: %s\nDIV: %04X TIMA: %02X\nTMA: %02X   TAC:  %02X",
		ebiten.ActualFPS(), ebiten.ActualTPS(),
		g.speed,
		registerState.AF, registerState.PC,
		registerState.BC, registerState.SP,
		registerState.DE,
		registerState.HL,
		registerState.IME,
		flagState,
		timerState.DIV,
		timerState.TIMA, timerState.TMA, timerState.TAC,
	)

	debugPanel.Fill(color.RGBA{53, 83, 10, 1})
	ebitenutil.DebugPrintAt(debugPanel, stats, 4, 4)

	drawOpt := &ebiten.DrawImageOptions{}
	drawOpt.GeoM.Translate(SCREEN_WIDTH, 0)
	screen.DrawImage(debugPanel, drawOpt)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	if g.debugOn {
		return SCREEN_WIDTH + DEBUG_PANEL_WIDTH, SCREEN_HEIGHT
	} else {
		return SCREEN_WIDTH, SCREEN_HEIGHT
	}
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

func (g *Game) updatePauseState() {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.pause = !g.pause
	}
}

func (g *Game) updateSpeed() {
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.speed = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.speed = 2
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.speed = 3
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.Key4) {
		g.speed = 4
	}
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.Key5) {
		g.speed = 5
	}
}

func (g *Game) updateFilter() {
	if ebiten.IsKeyPressed(ebiten.KeyControl) && inpututil.IsKeyJustPressed(ebiten.KeyG) {
		g.originalFilterEnabled = !g.originalFilterEnabled
	}
}
