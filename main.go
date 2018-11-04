package main

import (
	"log"

	"github.com/fatih/color"
)

// Delay - milliseconds to wait before begin next job after a file change
const Delay = 1000

func colorSuccess(format string) string {
	return color.GreenString(format)
}

func colorFail(format string) string {
	return color.RedString(format)
}

func main() {
	result := true

	if result != false {
		log.Println(colorSuccess("Success!"))
	} else {
		log.Println(colorFail("Fail!"))
	}
}
