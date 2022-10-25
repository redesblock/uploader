package watcher

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/redesblock/uploader/core/util"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

type HandleFunc func(event fsnotify.Event) error

type Watcher interface {
	Start(handlers ...HandleFunc) error
	Close() error
	AddRecursive(path string) error
	RemoveRecursive(path string) error
}

func New(ignoreHidden bool) Watcher {
	return &watch{
		ignoreHidden: ignoreHidden,
	}
}

type watch struct {
	*fsnotify.Watcher

	wg        sync.WaitGroup
	done      chan struct{}
	isRunning bool

	ignoreHidden bool                // ignore hidden files or not.
	ignored      map[string]struct{} // ignored files or directories.
}

func (w *watch) Start(handlers ...HandleFunc) error {
	if w.isRunning {
		return nil
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	w.Watcher = watcher
	w.done = make(chan struct{})
	w.isRunning = true

	w.wg.Add(1)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Debugf("watcher: add event %v", event)
				if event.Op&fsnotify.Create != 0 {
					if err := w.AddRecursive(event.Name); err != nil {
						log.Errorf("watcher: create event %v error %v", event, err)
					}
				}

				for _, handler := range handlers {
					if err := handler(event); err != nil {
						log.Errorf("watcher: handle event %v error %v", event, err)
					}
				}
			case err := <-watcher.Errors:
				log.Errorf("watcher: notify error %v", err)
			case <-w.done:
				w.Watcher.Close()
				return
			}
		}
	}()

	return nil
}

func (w *watch) Close() error {
	if !w.isRunning {
		return nil
	}

	w.isRunning = false
	close(w.done)
	w.wg.Wait()
	return nil
}

func (w *watch) AddRecursive(path string) error {
	if !w.isRunning {
		return fmt.Errorf("watcher instance not running, please start")
	}
	err := filepath.Walk(path, func(walkPath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			isHidden, err := util.IsHiddenFile(walkPath)
			if err != nil {
				return fmt.Errorf("isHidden error %v", err)
			}
			if w.ignoreHidden && isHidden {
				return filepath.SkipDir
			}

			if err = w.Add(walkPath); err != nil {
				return err
			}
			log.Infof("watcher: add path %s", walkPath)
		}
		return nil
	})
	return err
}

func (w *watch) RemoveRecursive(path string) error {
	if !w.isRunning {
		return fmt.Errorf("watcher instance not running, please start")
	}
	err := filepath.Walk(path, func(walkPath string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if fi.IsDir() {
			isHidden, err := util.IsHiddenFile(walkPath)
			if err != nil {
				return fmt.Errorf("isHidden error %v", err)
			}
			if w.ignoreHidden && isHidden {
				return filepath.SkipDir
			}

			if err = w.Remove(walkPath); err != nil {
				return err
			}
			log.Infof("watcher: remove path %s", walkPath)
		}
		return nil
	})
	return err
}
