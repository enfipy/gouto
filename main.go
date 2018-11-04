package main

import (
	"flag"
	"log"
	"os"

	"github.com/fatih/color"
)

// Delay - milliseconds to wait before begin next job after a file change
const Delay = 1000

// Pattern - watched files extensions pattern
const Pattern = `(.+\.go)$`

var (
	flagDirectory = flag.String("dir", ".", "Directory to watch for changes")
	flagOutput    = flag.String("out", ".", "Output directory for binary after build")
	flagRun       = flag.String("run", "", "Command to run after build")
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

	result := true

	if result != false {
		log.Println(colorSuccess("Success!"))
	} else {
		log.Println(colorFail("Fail!"))
	}
}
