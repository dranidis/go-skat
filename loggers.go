package main

import (
	"fmt"
	"github.com/fatih/color"
	"io"
	"os"
	"log"
)

var logFile io.Writer = nil

var debugTacticsLogFlag = false
var gameLogFlag = true
var fileLogFlag = true

var debugTacticsInMM = false


var debugTacticsLogFlagStored = false
var gameLogFlagStored = true
var fileLogFlagStored = true


var htmlLogFlag = false
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

func debugMinmaxLog(format string, a ...interface{}) {
	debugTacticsLog("MM: " + format, a...)
}

func createFile(logFileName string) *os.File {
	file, err := os.Create(logFileName)
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	return file
}


func disableLogs() {
	debugTacticsLogFlagStored = debugTacticsLogFlag
	gameLogFlagStored = gameLogFlag
	fileLogFlagStored = fileLogFlag
	if !debugTacticsInMM {
		debugTacticsLogFlag = false
		gameLogFlag = false
		fileLogFlag = false
	}
}

func restoreLogs() {
	debugTacticsLogFlag = debugTacticsLogFlagStored
	gameLogFlag = gameLogFlagStored
	fileLogFlag = fileLogFlagStored
}