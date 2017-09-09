package main

import (
	"fmt"
	"github.com/fatih/color"
	"io"
)

var logFile io.Writer = nil
var debugTacticsLogFlag = false
var gameLogFlag = true
var fileLogFlag = true
var htmlLogFlag = true
var logFileName = "gameLog.txt"

func logToFile(format string, a ...interface{}) {
	if fileLogFlag && logFile != nil {
		fmt.Fprintf(logFile, format, a...)
	}
}

func bidLog(format string, a ...interface{}) {
	if gameLogFlag {
		fmt.Printf(format, a...)
	}
	logToFile(format, a...)
}

func gameLog(format string, a ...interface{}) {
	if gameLogFlag {
		fmt.Printf(format, a...)
	}
	logToFile(format, a...)
}

func htmlLog(format string, a ...interface{}) {
	if htmlLogFlag {
		red := color.New(color.Bold, color.FgYellow).SprintFunc()
		s := fmt.Sprintf(format, a...)
		fmt.Printf(red(s))
	}
}

func debugTacticsLog(format string, a ...interface{}) {
	if debugTacticsLogFlag {
		fmt.Printf(format, a...)
	}
	logToFile(format, a...)
}

