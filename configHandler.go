package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/log"
	"github.com/ilyakaznacheev/cleanenv"
)

// Config represents the overall configuration structure
type Config struct {
	Delay    int `yaml:"delay" env-default:"10"`
	Watchers []struct {
		File FileConfig `yaml:"file"`
	} `yaml:"watchers"`
}

// FileConfig represents the configuration for each file entry
type FileConfig struct {
	Name  string   `yaml:"name"`
	Exact bool     `yaml:"exact" env-default:"false"`
	Cmds  []string `yaml:"cmds"`
}

var (
	cfg Config

	input string

	targetingStatus = make(chan bool)
	watcherStatus   = make(chan bool)
)

var defaultConfig = `# K-9 config

delay: 10

watchers:
- file:
    name: main.go
    cmds:
    - cmd /c echo Hello from K-9!

`

func parseConfig() {
	err := cleanenv.ReadConfig("k-9.yml", &cfg)
	if err != nil {
		err = cleanenv.ReadConfig("k-9.yaml", &cfg)
		if err != nil {
			k9p.Error("Seems like K-9 can't find or read a config file in this directory. Would you like to initialize one? [y/yes] or [n/no]")
			fmt.Scanln(&input)
			if input == "y" || input == "yes" {
				initConfig()
				k9p.Error("Edit the k-9.yml file and run the program again. Exiting...")
				os.Exit(1)
			} else {
				log.Fatal("Exiting...")
			}
		}
	}

	k9p.Info("⚙️  K-9 config loaded")
	// log.Warn(cfg.Watchers)

	go initTargetStatus()
	go initWatcherStatus()

	initializeWatchers()
	// for _, watcher := range cfg.Watchers {
	// 	log.Info(watcher.File.Type)
	// 	log.Info(watcher.File.Cmds)
	// }
}

// status logging
func initTargetStatus() {

	for i := 0; i < len(cfg.Watchers); i++ {
		<-targetingStatus
	}
	k9p.Info("📂 Watcher targets found")
}

func initWatcherStatus() {

	for i := 0; i < len(cfg.Watchers); i++ {
		<-watcherStatus
	}
	k9p.Info("🔁 All watchers are online")
}

// config initialization
func initConfig() {
	f, err := os.Create("k-9.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Write([]byte(defaultConfig))
}
