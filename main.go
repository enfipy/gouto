package main

import (
	"bufio"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/fsnotify/fsnotify"
)

// Delay - milliseconds to wait before restart after a file change
const Delay = 100

var (
	flagDirectory = flag.String("dir", "./", "Directory to watch for changes")
	flagOutput    = flag.String("out", "./cmd/app", "Output directory for binary after build")
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

	setupWatcher()
}

func setupWatcher() {
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

	watchFiles(watcher)

	var currentProcess *os.Process
	if build() {
		currentProcess = runBinary().Process
	}

	printSuccess("Waiting for changes...")
	for {
		select {
		case err := <-watcher.Errors:
			printFail("Failed to reload: ", err.Error())
		case event := <-watcher.Events:
			if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
				time.Sleep(Delay * time.Millisecond)

				if currentProcess != nil {
					killProcess(currentProcess)
				}

				printSuccess("Restarting...")
				if build() {
					currentProcess = runBinary().Process
					printSuccess("Started")
				}
			}
		}
	}
}

func watchFiles(watcher *fsnotify.Watcher) {
	err := filepath.Walk(*flagDirectory, func(path string, info os.FileInfo, err error) error {
		if err == nil && info.IsDir() {
			return watcher.Add(path)
		}
		return err
	})
	if err != nil {
		printFail("Failed to add inner folders for watching: ", err.Error())
		os.Exit(1)
	}

	err = watcher.Add(*flagDirectory)
	if err != nil {
		printFail("Failed to add folder for watching: ", err.Error())
		os.Exit(1)
	}
}

func build() bool {
	args := strings.Split(*flagBuild, " ")
	if len(args) == 0 {
		return true
	}
	args = append(args, "-o", outDir)

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = *flagDirectory

	if output, err := cmd.CombinedOutput(); err != nil {
		printFail("Failed while building: ", string(output))
		return false
	}
	return true
}

func runBinary() *exec.Cmd {
	var stdout, stderr io.ReadCloser
	var err error

	cmd := exec.Command(outDir)

	if stdout, err = cmd.StdoutPipe(); err != nil {
		printFail("Failed get stdout pipe: ", err.Error())
		return nil
	}
	if stderr, err = cmd.StderrPipe(); err != nil {
		printFail("Failed get stderr pipe: ", err.Error())
		return nil
	}

	go logger(stdout)
	go logger(stderr)

	if err = cmd.Start(); err != nil {
		printFail("Failed while running: ", err.Error())
		return nil
	}
	return cmd
}

func logger(pipe io.ReadCloser) {
	reader := bufio.NewReader(pipe)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		log.Print(line)
	}
}

func killProcess(process *os.Process) {
	if err := process.Kill(); err != nil {
		printFail("Can not kill current process: ", err.Error())
	}
	if _, err := process.Wait(); err != nil {
		printFail("Can not wait for exiting of current process")
	}
}
