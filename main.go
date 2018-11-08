package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

// Delay - milliseconds to wait before begin next job after a file change
const Delay = 1000

// Pattern - watched files extensions pattern
const Pattern = `(.+\.go)$`

var (
	flagDirectory = flag.String("dir", ".", "Directory to watch for changes")
	flagStart     = flag.String("start", "", "Command to run after build")
	flagBuild     = flag.String("build", "go build", "Command to rebuild after changes")
)

func colorSuccess(format string) string {
	return color.GreenString(format)
}

func colorFail(format string) string {
	return color.RedString(format)
}

func main() {
	flag.Parse()

	if *flagDirectory == "" {
		log.Println(colorFail("Directory flag is required"))
		os.Exit(1)
	}

	watchFiles()
}

func watchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(colorFail("Failed to setup watcher"))
		os.Exit(1)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case err := <-watcher.Errors:
				log.Println(colorFail("Failed to reload"), err)
			case event := <-watcher.Events:
				onChange(event)
			}
		}
	}()

	if err := watcher.Add(*flagDirectory); err != nil {
		log.Println(colorFail("Failed to add provided folder to watcher"))
		os.Exit(1)
	}

	log.Println(colorSuccess("Waiting for changes..."))
	<-done
}

func onChange(event fsnotify.Event) {
	trigger := event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create

	if trigger {
		time.Sleep(100 * time.Millisecond)
		build()
		start()
	}
}

func build() {
	log.Println(colorSuccess("Success build"))
}

func start() {
	log.Println(colorSuccess("Success start"))
}
