package main

import (
	"github.com/qhdwight/biomequest/game"
	"runtime"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	game.Start()
}
