package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	filename *string
	cmds     *[]string
}

// var executionCounter = 0

func (w *Watcher) findFilePath() (string, error) {
	files, err := os.ReadDir(".")
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if *w.filename != "" && strings.Contains(file.Name(), *w.filename) {
			// log.Info("📂 watcher target: " + file.Name() + "found")
			return filepath.Join(".", file.Name()), nil
		}
	}

	return "", nil
}

func (w *Watcher) executeCommands() {
	for _, cmd := range *w.cmds {
		// log.Infof("Executing command: %s for file type: %s\n", cmd, *w.filetype)
		// Here you can execute the command, e.g., using os/exec package
		command := strings.Split(cmd, " ")

		cmdOutput, err := exec.Command(command[0], command[1:]...).Output()
		if err != nil {
			log.Warn(err)
		} else {
			log.Info(*w.filename + " changed" + "\n" + string(cmdOutput))
		}
	}
}

// init initializes the Watcher by finding the file path, creating a new fsnotify watcher,
// adding the file path to the watcher, and starting an infinite loop to monitor file events.
//
// Returns an error if there was an issue finding the file path, creating the watcher, or adding the file path to the watcher.
func (w *Watcher) init() error {
	filePath, err := w.findFilePath()
	if err != nil {
		return err
	}

	if filePath == "" {
		log.Fatal("Specified file not found")
	}

	log.Info("📂 watcher target: " + *w.filename + "found")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	err = watcher.Add(filePath)
	if err != nil {
		return err
	}

	log.Info("🔁 watcher created and started")

	var lastEventTime time.Time

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return nil
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if time.Since(lastEventTime) > time.Duration(cfg.Delay)*time.Second {
					lastEventTime = time.Now()
					// log.Info("modified file: %s\n", event.Name)
					w.executeCommands()
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return nil
			}
			log.Warn(err)
		}
	}
}

func initializeWatchers() {
	for _, watcher := range cfg.Watchers {
		w := &Watcher{
			filename: &watcher.File.Name,
			cmds:     &watcher.File.Cmds,
		}
		go func(w *Watcher) {
			if err := w.init(); err != nil {
				log.Fatalf("Error initializing watcher: %v\n", err)
			}
		}(w)
	}

	log.Info("👁️  Watchers initialized")

	// Block main goroutine forever.
	select {}
}
