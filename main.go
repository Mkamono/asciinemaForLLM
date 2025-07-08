package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Mkamono/asciinemaForLLM/internal/cmd"
)

func main() {
	args := os.Args[1:] // Remove program name

	// Handle no arguments (default format from stdin)
	if len(args) == 0 {
		if err := cmd.RunFormat("structured"); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Parse arguments
	command := args[0]
	cleanup := false
	outputFormat := "structured"

	// Parse flags
	i := 0
	for i < len(args) {
		arg := args[i]
		if arg == "--cleanup" {
			cleanup = true
			args = append(args[:i], args[i+1:]...)
			continue
		} else if strings.HasPrefix(arg, "--output=") {
			outputFormat = strings.TrimPrefix(arg, "--output=")
			args = append(args[:i], args[i+1:]...)
			continue
		}
		i++
	}

	switch command {
	case "format":
		if err := cmd.RunFormat(outputFormat); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "record":
		var outputFile string
		if len(args) > 1 {
			outputFile = args[1]
		}
		if err := cmd.RunRecord(outputFile, cleanup, outputFormat); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "file":
		if len(args) < 2 {
			fmt.Fprintf(os.Stderr, "Error: file command requires input file\n")
			cmd.ShowUsage()
			os.Exit(1)
		}
		inputFile := args[1]
		var outputFile string
		if len(args) > 2 && !strings.HasPrefix(args[2], "--") {
			outputFile = args[2]
		}
		if err := cmd.RunFormatFile(inputFile, outputFile, cleanup, outputFormat); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

	case "-h", "--help", "help":
		cmd.ShowUsage()

	default:
		fmt.Fprintf(os.Stderr, "Error: unknown command '%s'\n", command)
		cmd.ShowUsage()
		os.Exit(1)
	}
}