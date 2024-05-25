package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/ilyakaznacheev/cleanenv"
)

// FileConfig represents the configuration for each file entry
type FileConfig struct {
	Name string   `yaml:"name"`
	Cmds []string `yaml:"cmds"`
}

// Config represents the overall configuration structure
type Config struct {
	Watchers []struct {
		File FileConfig `yaml:"file"`
	} `yaml:"watchers"`
}

var cfg Config
var input string

var defaultConfig = `# K-9 config

watchers:
- file:
    name: main.go
    cmds:
    - go test

`

func parseConfig() {
	err := cleanenv.ReadConfig("k-9.yml", &cfg)
	if err != nil {
		log.Info("Seems like K-9 can't find or read a config file in this directory. Would you like to initialize one? [y/yes] or [n/no]")
		fmt.Scanln(&input)
		if input == "y" || input == "yes" {
			initConfig()
			log.Info("Edit the k-9.yml file and run the program again. Exiting...")
			os.Exit(1)
		} else {
			log.Fatal("Exiting...")
		}
	}

	initializeWatchers()

	// for _, watcher := range cfg.Watchers {
	// 	log.Info(watcher.File.Type)
	// 	log.Info(watcher.File.Cmds)
	// }
}

func initConfig() {
	f, err := os.Create("k-9.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Write([]byte(defaultConfig))
}
