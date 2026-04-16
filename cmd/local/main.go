package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"goboy/app"
	"goboy/app/cartridge"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// Parse flags
	scale := flag.Int("scale", 5, "Multiplier for screen size")
	debug := flag.Bool("debug", false, "Prints debug info to screen")
	noSync := flag.Bool("nosync", false, "Will draw to screen even if frame not fully rendered (possible screen tearing)")
	saveFile := flag.String("save-file", "", "File name of save file")
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
		os.Exit(1)
	}

	print(cart.String())

	game := app.NewGame(cart, *debug, !*noSync, *scale)

	if saveFile != nil {
		saveState, err := app.LoadSaveFromFile(*saveFile)
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}

		game.LoadSaveState(saveState)
	}

	ebiten.SetWindowTitle("GoBoy")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
