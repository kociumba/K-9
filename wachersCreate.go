package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	filename *string
	cmds     *[]string
}

// var executionCounter = 0
var shouldDetect []string

func (w *Watcher) findFilePath() ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	files, err := os.ReadDir(wd)
	if err != nil {
		return nil, err
	}

	var matchingFiles []string

	// This deep fried piece of shit loop holds this together 🥲
	for _, file := range files {
		if *w.filename != "" {
			dotIndex := strings.Index(*w.filename, ".")

			if dotIndex == -1 {
				// by name only
				if strings.Contains(file.Name(), *w.filename) {
					matchingFiles = append(matchingFiles, file.Name())
				}
			} else if dotIndex == 0 {
				// by file extension
				if strings.Contains(file.Name(), *w.filename) {
					matchingFiles = append(matchingFiles, file.Name())
				}
			} else {
				// by full file name
				if strings.Contains(file.Name(), *w.filename) {
					// returns first file matching the name
					// log.Warn(file.Name())
					shouldDetect = append(shouldDetect, file.Name())
					return []string{filepath.Join(".", file.Name())}, nil
				}
			}
		}
	}

	if len(matchingFiles) > 0 {
		// avoid createing multiple watchers for the same file
		matchingFilesTemp := make([]string, 0, len(matchingFiles))
		for _, file := range matchingFiles {
			if !slices.Contains(shouldDetect, file) {
				matchingFilesTemp = append(matchingFilesTemp, file)
			}
		}
		matchingFiles = matchingFilesTemp
		shouldDetect = append(shouldDetect, matchingFiles...)
		// log.Warn(matchingFiles)
		return matchingFiles, nil
	}

	return nil, nil
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

	if filePath == nil || filePath[0] == "" {
		log.Fatal("Specified file/s not found")
	}

	targetingStatus <- true

	// log.Info("📂 watcher target: " + *w.filename + " found")

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}
	defer watcher.Close()

	for i := 0; i < len(filePath); i++ {
		err = watcher.Add(filePath[i])
		if err != nil {
			return err
		}
	}

	watcherStatus <- true

	// log.Info("🔁 watcher created and started")

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
