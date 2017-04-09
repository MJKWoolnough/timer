package main

import (
	"errors"
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/screensaver"
	"github.com/MJKWoolnough/engine"
	//_ "github.com/MJKWoolnough/engine/graphics/sdl1"
	_ "github.com/MJKWoolnough/engine/windows/sdl1"
)

func main() {
	err := run()
	if err != nil {
		log.Println(err)
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
	if err := engine.Init(engine.Config{
		Monitor: monitors[0],
		Mode:    modes[len(modes)-1],
		Title:   "Timer",
	}); err != nil {
		return err
	}
	conn, err := xgb.NewConn()
	if err != nil {
		return err
	}
	screensaver.Init(conn)
	screensaver.Suspend(conn, true)
	setup()
	engine.Loop(loop)
	screensaver.Suspend(conn, false)
	conn.Close()
	return engine.Uninit()
}

var last = -1

func loop(w, h int, t float64) bool {
	if engine.KeyPressed(engine.KeyEscape) {
		engine.Close()
		return false
	}
	s := int(t)
	if s != last {
		render(w, h, s)
		last = s
	}
	return true
}
