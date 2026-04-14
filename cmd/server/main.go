package main

import (
	_ "embed"
	"log"

	"goboy/app"
	"goboy/app/cartridge"

	"github.com/hajimehoshi/ebiten/v2"
)

//go:embed "roms/wordyl.gb"
var romData []byte

func main() {
	cart := cartridge.NewCartridge("Wordyl", romData)

	game := app.NewGame(cart, true, true, 3)

	print(cart.String())
	ebiten.SetWindowSize(160*2, 144*2)
	ebiten.SetWindowTitle("GoBoy")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
