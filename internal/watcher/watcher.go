package watcher

import (
	"fmt"
	"goTSV/config"
	"os"
	"strings"
	"sync"
	"time"
)

func NewWatcher(c config.Config) *Watcher {
	return &Watcher{timer: c.RefreshInterval, fromDir: c.DirectoryFrom, files: make(map[string]bool)}
}

type Watcher struct {
	mutex   sync.RWMutex
	timer   int
	fromDir string
	files   map[string]bool
}

func (s *Watcher) Scan(out chan string) {
	timer := time.NewTicker(time.Duration(s.timer) * time.Second)

	defer timer.Stop()

	for range timer.C {
		filesFromDir, err := os.ReadDir(s.fromDir)
		if err != nil {
			fmt.Errorf("error reading %w", err)
			return
		}
		for _, file := range filesFromDir {
			if strings.HasSuffix(file.Name(), ".tsv") && !file.IsDir() {
				s.mutex.Lock()
				_, ok := s.files[file.Name()]
				if ok {
					s.mutex.Unlock()
					continue
				} else {
					s.files[file.Name()] = true
					s.mutex.Unlock()
				}
				select {
				case out <- file.Name():
				}
			}
		}
	}
}
