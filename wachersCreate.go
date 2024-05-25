package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	filetype *string
	cmds     *[]string
}

var executionCounter = 0

func (w *Watcher) findFilePath() (string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if *w.filetype != "" && strings.Contains(file.Name(), *w.filetype) {
			return filepath.Join(".", file.Name()), nil
		}
	}

	return "", nil
}

func (w *Watcher) executeCommands() {
	for _, cmd := range *w.cmds {
		log.Infof("Executing command: %s for file type: %s\n", cmd, *w.filetype)
		// Here you can execute the command, e.g., using os/exec package
		command := strings.Split(cmd, " ")

		cmdOutput, err := exec.Command(command[0], command[1:]...).Output()
		if err != nil {
			log.Warn(err)
		} else {
			log.Info("\n" + string(cmdOutput))
		}
	}
}

func (w *Watcher) init() error {
	filePath, err := w.findFilePath()
	if err != nil {
		return err
	}

	if filePath == "" {
		log.Fatal("No file with specified type found in the working dir")
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(filePath)
	if err != nil {
		return err
	}

	var lastEventTime time.Time

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if time.Since(lastEventTime) > 1*time.Second {
					lastEventTime = time.Now()
					log.Infof("INFO modified file: %s\n", event.Name)
					w.executeCommands()
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Infof("error: %v\n", err)
		}
	}
}

func initializeWatchers() {
	for _, watcher := range cfg.Watchers {
		w := &Watcher{
			filetype: &watcher.File.Name,
			cmds:     &watcher.File.Cmds,
		}
		go func(w *Watcher) {
			if err := w.init(); err != nil {
				log.Fatalf("Error initializing watcher: %v\n", err)
			}
		}(w)
	}

	// Block main goroutine forever.
	select {}
}
