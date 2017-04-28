package main

import (
	"pkg.deepin.io/lib/log"
	"simple-session/display"
	"simple-session/keybinding"
)

var (
	logger = log.NewLogger("SimpleSession")
)

func main() {
	err := keybinding.Load(logger)
	if err != nil {
		logger.Error("Failed to load keybinding:", err)
	}
	err = display.Load(logger)
	if err != nil {
		logger.Error("Failed to load display:", err)
	}
}

func doToggleDebug() {
	if logger.GetLogLevel() == log.LevelDebug {
		logger.SetLogLevel(log.LevelInfo)
	} else {
		logger.SetLogLevel(log.LevelDebug)
	}
}

func getOutputInfos() string {
	return display.GetOutputInfos()
}
