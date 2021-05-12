package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"syscall"
)

//ExecShell is basically equivalent to "echo "$stdinPipe" | command arg"
func ExecShell(command string, arg string, stdinPipe string) (string, int) {
	var exitCode int

	cmd := exec.Command(command, arg)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, stdinPipe)
	}()

	output, err := cmd.CombinedOutput()
	if err != nil {
		// if failed because of exit status, don't log to console
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.Sys().(syscall.WaitStatus).ExitStatus()
		} else {
			fmt.Printf("err is %v\n", err)
		}
	} else {
		exitCode = cmd.ProcessState.Sys().(syscall.WaitStatus).ExitStatus()
	}

	return string(output), exitCode
}

// Opens text files and returns its content as string
// Also is able to return os.Stdin if argument "stdin" is passed
func ReadTextFile(filename string) string {
	var fileContent string

	if filename == "stdin" {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			fileContent += scanner.Text() + "\n"
		}
	} else {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatal(err)
		}
		fileContent = string(data)
	}

	return fileContent
}
