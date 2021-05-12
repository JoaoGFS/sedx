package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Insuficient args")
	}

	scriptFile := os.Args[1]
	inputText, _ := ExecShell("sed", "s/^[ \t]*//g;s/^$//g", ReadTextFile(os.Args[2]))
	lineNumber := 1

	var ifConditionFlag bool
	var insideIfFlag bool
	outputText := inputText

	f, _ := os.Open(scriptFile)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		currentLine := scanner.Text()

		if strings.HasPrefix(currentLine, "if") {
			insideIfFlag = true
			_, exitCode := ExecShell("grep", strings.Split(currentLine, " ")[1], inputText)
			if exitCode == 0 {
				ifConditionFlag = true
			} else {
				ifConditionFlag = false
			}
			// fmt.Printf("condition %s was %v\n", strings.Split(currentLine, " ")[1], ifConditionFlag)
		} else if strings.HasPrefix(currentLine, "end") {
			insideIfFlag = false
		} else if insideIfFlag && ifConditionFlag {
			outputText, _ = ExecShell("sed", currentLine, outputText)
		} else if !insideIfFlag {
			outputText, _ = ExecShell("sed", currentLine, outputText)
		}

		lineNumber++

	}

	fmt.Println(outputText)

}
