package main

import (
	"goTSV/config"
	"goTSV/internal/handler"
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
	s := service.NewService(st, w, cnfg)
	h := handler.NewHandler(s, cnfg)
	go func() {
		err := s.Scanner()
		if err != nil {
			return
		}
	}()
	h.Start()
}
