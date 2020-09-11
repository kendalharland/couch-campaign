package main

import (
	"couchcampaign"
	"encoding/json"
	"errors"
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/font"
)

const (
	title          = "Couch Campaign"
	version        = 1.0
	configFilename = "couchcampaign.json"
)

type Button = pixelgl.Button

const (
	buttonAccept  Button = pixelgl.KeyD
	buttonReject  Button = pixelgl.KeyA
	buttonNewGame Button = pixelgl.KeyN
	buttonQuit    Button = pixelgl.KeyQ
)

var defaultConfig = couchcampaign.Configuration{
	WindowWidth:  500,
	WindowHeight: 650,
}

var gameFont font.Face

func init() {
	f, err := loadTTF("assets/xolonium/ttf/Xolonium-Regular.ttf", 12)
	if err != nil {
		panic(err)
	}
	gameFont = f
}

func main() {
	pixelgl.Run(func() {
		if err := run(); err != nil {
			log.Fatal(err)
		}
	})
}

func run() error {
	config, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}
	win, err := pixelgl.NewWindow(pixelgl.WindowConfig{
		Title:  fmt.Sprintf("%s v%v", title, version),
		Bounds: pixel.R(0, 0, config.WindowWidth, config.WindowHeight),
	})
	if err != nil {
		return err
	}

	game, err := couchcampaign.NewGame()
	if err != nil {
		return err
	}

	fps := time.Tick(time.Second / 120)

	var state gameState = gameMainMenuState

	pic, err := loadPicture("assets/droid-zapper/wake.png")
	if err != nil {
		panic(err)
	}
	sprite, err := NewSpriteFromMeasurements(pic, DroidZapperWakeMeasurements, 0)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		<-fps

		win.Clear(color.Black)
		next, err := state(win, game)
		if err != nil {
			return err
		}
		center := win.Bounds().Center()
		sprite.Draw(win, pixel.IM.Moved(center).Scaled(center, 2))
		state = next
		win.Update()
	}
	return nil
}

type gameState func(*pixelgl.Window, *couchcampaign.Game) (gameState, error)

func gameMainMenuState(win *pixelgl.Window, game *couchcampaign.Game) (gameState, error) {
	canvas := text.New(pixel.V(100, 500), text.NewAtlas(gameFont, text.ASCII))
	fmt.Fprintln(canvas, "Couch Campaign")
	fmt.Fprintln(canvas, "")
	fmt.Fprintf(canvas, "[%s] New game\n", buttonNewGame.String())
	fmt.Fprintf(canvas, "[%s] Quit\n", buttonQuit.String())
	canvas.Draw(win, pixel.IM)

	switch {
	case win.JustPressed(buttonNewGame):
		return gamePlayingState, nil
	case win.JustPressed(buttonQuit):
		return nil, errors.New("exit")
	default:
		return gameMainMenuState, nil
	}
}

func gamePlayingState(win *pixelgl.Window, game *couchcampaign.Game) (gameState, error) {
	getInput := func(win *pixelgl.Window) couchcampaign.Input {
		if win.JustPressed(buttonAccept) {
			return couchcampaign.InputCardAccepted
		}
		if win.JustPressed(buttonReject) {
			return couchcampaign.InputCardRejected
		}
		return couchcampaign.NoInput
	}

	state := game.GetState()

	progress := text.New(pixel.V(100, 500), text.NewAtlas(gameFont, text.ASCII))
	fmt.Fprintf(progress, "Character: %v\n", state.Character)
	fmt.Fprintf(progress, "Days survived: %v\n", state.CharacterLifespan)
	progress.Draw(win, pixel.IM)

	input := getInput(win)
	if input == couchcampaign.NoInput {
		return gamePlayingState, nil
	}
	if err := game.HandleInput(input); err != nil {
		return gamePlayingState, err
	}
	if state.IsFailed() {
		return gameFailedState, nil
	}
	return gamePlayingState, nil
}

func gameFailedState(win *pixelgl.Window, game *couchcampaign.Game) (gameState, error) {
	basicAtlas := text.NewAtlas(gameFont, text.ASCII)
	basicTxt := text.New(pixel.V(100, 500), basicAtlas)
	fmt.Fprintln(basicTxt, "you failed")
	basicTxt.Draw(win, pixel.IM)

	return gameFailedState, nil
}

func loadConfig() (*couchcampaign.Configuration, error) {
	cp, err := configPath()
	if err != nil {
		return nil, err
	}

	if !fileExists(cp) {
		if err := writeDefaultConfig(cp); err != nil {
			return nil, err
		}
	}
	data, err := ioutil.ReadFile(cp)
	if err != nil {
		return nil, fmt.Errorf("read %s: %v", cp, err)
	}
	config := new(couchcampaign.Configuration)
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("unmarshal %s: %v", cp, err)
	}
	return config, nil
}

func writeDefaultConfig(filename string) error {
	data, err := json.Marshal(defaultConfig)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, data, 0755); err != nil {
		return err
	}
	return nil
}

func configPath() (string, error) {
	dir, err := os.UserHomeDir()
	if err != nil || dir == "" {
		dir = "./"
	}
	return filepath.Abs(filepath.Join(dir, configFilename))
}

func fileExists(name string) bool {
	stat, err := os.Stat(name)
	if err != nil {
		if !os.IsNotExist(err) {
			log.Printf("warning: could not determine if path %q exists: %v", name, err)
		}
		return false
	}
	return !stat.IsDir()
}
