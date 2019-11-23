package main

import (
	"github.com/qhdwight/go-game/game"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	game.Start()
}
