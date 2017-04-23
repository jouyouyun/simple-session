package main

import (
	"pkg.deepin.io/lib/log"
	"simple-session/keybinding"
)

var (
	logger = log.NewLogger("SimpleSession")
)

func main() {
	keybinding.Load(logger)
}

func doToggleDebug() {
	if logger.GetLogLevel() == log.LevelDebug {
		logger.SetLogLevel(log.LevelInfo)
	} else {
		logger.SetLogLevel(log.LevelDebug)
	}
}
