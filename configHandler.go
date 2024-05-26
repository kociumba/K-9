package main

import (
	"fmt"
	"os"

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
	Name string   `yaml:"name"`
	Cmds []string `yaml:"cmds"`
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

	log.Info("⚙️  K-9 config loaded")

	go initTargetStatus()
	go initWatcherStatus()

	initializeWatchers()
	// for _, watcher := range cfg.Watchers {
	// 	log.Info(watcher.File.Type)
	// 	log.Info(watcher.File.Cmds)
	// }
}

func initTargetStatus() {

	for i := 0; i < len(cfg.Watchers); i++ {
		<-targetingStatus
	}
	log.Info("📂 Watcher targets found")
}

func initWatcherStatus() {

	for i := 0; i < len(cfg.Watchers); i++ {
		<-watcherStatus
	}
	log.Info("🔁 All watchers are online")
}

func initConfig() {
	f, err := os.Create("k-9.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	f.Write([]byte(defaultConfig))
}
