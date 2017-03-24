package main

import (
	"errors"
	"log"

	"github.com/MJKWoolnough/engine"
	_ "github.com/MJKWoolnough/engine/graphics/gles2"
	_ "github.com/MJKWoolnough/engine/windows/glfw32"
)

func main() {
	err := run()
	if err != nil {
		log.Printf(err)
	}
}

func run() error {
	monitors := engine.GetMonitors()
	if len(monitors) == 0 {
		return errors.New("no monitor")
	}
	modes := monitors[0].GetModes()
	if len(modes) == 0 {
		return errors.New("no modes")
	}
	return engine.Loop(engine.Config{
		Monitor: monitors[0],
		Mode:    modes[len(modes)-1],
		Title:   "Timer",
	}, loop)
}

func loop(w, h int, t float64) {
}
