package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

// Delay - milliseconds to wait before begin next job after a file change
const Delay = 1000

// Pattern - watched files extensions pattern
const Pattern = `(.+\.go)$`

var (
	flagDirectory = flag.String("dir", "./", "Directory to watch for changes")
	flagOutput    = flag.String("out", "./cmd/app", "Output directory for binary after build")
	flagArguments = flag.String("args", "", "Arguments to run binary after build")
	flagBuild     = flag.String("build", "go build", "Command to rebuild after changes")

	outDir = ""
)

func printSuccess(format string) {
	log.Println(color.GreenString(format))
}

func printFail(format string, errors ...string) {
	log.Println(color.RedString(format, errors))
}

func main() {
	flag.Parse()

	if *flagDirectory == "" {
		printFail("Directory flag is required")
		os.Exit(1)
	}

	watchFiles()
}

func watchFiles() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		printFail("Failed to setup watcher: ", err.Error())
		os.Exit(1)
	}
	defer watcher.Close()

	outDir, err = filepath.Abs(*flagOutput)
	if outDir == "" {
		printFail("Invalid output directory: ", err.Error())
		os.Exit(1)
	}

	done := make(chan bool)
	go func() {
		for {
			select {
			case err := <-watcher.Errors:
				printFail("Failed to reload: ", err.Error())
			case event := <-watcher.Events:
				onChange(event)
			}
		}
	}()

	if err := watcher.Add(*flagDirectory); err != nil {
		printFail("Failed to add provided folder to watcher: ", err.Error())
		os.Exit(1)
	}

	if build() {
		start()
	}

	printSuccess("Waiting for changes...")
	<-done
}

func onChange(event fsnotify.Event) {
	trigger := event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create

	if trigger {
		time.Sleep(100 * time.Millisecond)
		printSuccess("Restarting...")
		if build() {
			start()
		}
	}
}

func build() bool {
	args := strings.Split(*flagBuild, " ")
	if len(args) == 0 {
		return true
	}
	// args = append(args, "-a")
	args = append(args, "-o", outDir)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = *flagDirectory

	if output, err := cmd.CombinedOutput(); err != nil {
		printFail("Failed while building: ", string(output))
		return false
	}
	return true
}

func start() {
	cmd := exec.Command(outDir)

	if err := cmd.Start(); err != nil {
		printFail("Failed while running: ", err.Error())
	}

	printSuccess("Started")

	// Todo: Make logging to console
}
