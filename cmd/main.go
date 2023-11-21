package main

import (
	"goTSV/config"
	"goTSV/internal/service"
	"goTSV/internal/storage"
	"goTSV/internal/watcher"
)

func main() {
	cnfg := config.New()
	st, err := storage.NewDBStorage(cnfg)
	if err != nil {
		return
	}
	w := watcher.NewWatcher(cnfg)
	s := service.NewService(*st, w, cnfg)
	s.Scanner()

}
