package main

import (
	"goTSV/config"
	"goTSV/internal/watcher"
)

func main() {
	cnfg := config.New()
	w := watcher.NewWatcher(cnfg)

}
