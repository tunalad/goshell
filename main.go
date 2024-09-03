package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var currentDir string

func init() {
	currentDir, _ = os.Getwd()
}

func shellLoop() {
	var line string
	var args []string
	var status int

	exit := false

	for !exit {
		fmt.Printf("goshell > ")

		line = readLine()
		args = splitLine(line)
		status = execLine(args)

		if status == -1 {
			exit = true
		}
	}
}

func readLine() string {
	var line string
	reader := bufio.NewReader(os.Stdin)

	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error ocurred on readLine")
		return ""
	}

	return line
}

func splitLine(line string) []string {
	line = strings.TrimSpace(line)

	var args []string
	var arg bytes.Buffer

	inQuote := false

	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '"':
			inQuote = !inQuote
		case ' ':
			if !inQuote && arg.Len() > 0 {
				args = append(args, arg.String())
				arg.Reset()
			}
		default:
			arg.WriteByte(line[i])
		}
	}

	if arg.Len() > 0 {
		args = append(args, arg.String())
	}

	return args
}

func execLine(args []string) int {
	status, output, err := launchProcess(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return 1
	}

	fmt.Fprint(os.Stdout, output)
	return status
}

func launchProcess(args []string) (int, string, error) {
	if len(args) == 0 {
		return 0, "", fmt.Errorf("no command provided")
	}

	biCode, biOutput, err := builtinCommands(args)
	if err != nil {
		if biCode == 123123123 {
			return 0, biOutput, nil
		}
	}

	if biCode == 0 && biOutput == "" && err == nil {
		return 0, "", nil
	}

	cmdName := args[0]

	cmdPath, err := exec.LookPath(cmdName)
	if err != nil {
		return 0, "", fmt.Errorf("command '%s' not found", cmdName)
	}

	cmd := exec.Command(cmdPath, args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode(), string(output), nil
		}
		return 0, "", fmt.Errorf("command execution failed: %v", err)
	}

	return 0, string(output), nil
}

func builtinCommands(args []string) (int, string, error) {
	switch args[0] {
	case "exit":
		os.Exit(0)
		return -1, "", nil
	case "cd":
		if len(args) > 1 {
			var err error
			currentDir, err = filepath.Abs(filepath.Join(currentDir, args[1]))
			if err != nil {
				return 123123123, "", err
			}
			if err := os.Chdir(currentDir); err != nil {
				return 123123123, "", err
			}
		} else {
			return 123123123, "", fmt.Errorf("cd: missing argument")
		}
	default:
		return 0, "", fmt.Errorf("not a default command")
	}
	return 0, "", nil
}

func main() {
	// load config files

	// loop
	shellLoop()

	// perform shotdown/cleanup shyt
	fmt.Println("exiting")
}
