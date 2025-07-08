package formatter

import (
	"fmt"
	"strings"
	"time"

	"github.com/Mkamono/asciinemaForLLM/internal/parser"
)

func FormatSession(session *parser.CastSession) {
	FormatSessionAsStructured(session)
}

func FormatSessionAsStructured(session *parser.CastSession) {
	fmt.Printf("Terminal Session (%s shell, %dx%d)\n", 
		extractShellName(session.Header.Env["SHELL"]), 
		session.Header.Width, 
		session.Header.Height)
	fmt.Printf("Recorded: %s\n", time.Unix(session.Header.Timestamp, 0).Format("2006-01-02 15:04:05"))
	if session.WorkingDir != "" {
		fmt.Printf("Working Directory: %s\n", session.WorkingDir)
	}
	fmt.Printf("\n")

	for _, cmd := range session.Commands {
		fmt.Printf("COMMAND: %s\n", cmd.Command)
		fmt.Printf("START TIME: %.3fs\n", cmd.StartTime)
		fmt.Printf("DURATION: %.3fs\n", cmd.EndTime-cmd.StartTime)
		fmt.Printf("EXIT CODE: %d\n", cmd.ExitCode)
		fmt.Printf("OUTPUT: ")
		if cmd.Output != "" {
			fmt.Printf("%s\n", cmd.Output)
		} else {
			fmt.Printf("(no output)\n")
		}
		fmt.Printf("\n")
	}
}

func FormatSessionAsCSV(session *parser.CastSession) {
	shell := extractShellName(session.Header.Env["SHELL"])
	width := session.Header.Width
	height := session.Header.Height
	recorded := time.Unix(session.Header.Timestamp, 0).Format("2006-01-02 15:04:05")
	workingDir := session.WorkingDir
	if workingDir == "" {
		workingDir = "(unknown)"
	}
	
	// CSV header
	fmt.Printf("shell,width,height,recorded,working_dir,command,start_time,duration,exit_code,output\n")
	
	// CSV rows
	for _, cmd := range session.Commands {
		output := cmd.Output
		if output == "" {
			output = "(no output)"
		}
		
		// Escape CSV fields properly
		fmt.Printf("%s,%d,%d,%s,%s,%s,%.3f,%.3f,%d,%s\n",
			csvEscape(shell),
			width,
			height,
			csvEscape(recorded),
			csvEscape(workingDir),
			csvEscape(cmd.Command),
			cmd.StartTime,
			cmd.EndTime-cmd.StartTime,
			cmd.ExitCode,
			csvEscape(output))
	}
}

func csvEscape(field string) string {
	// If field contains comma, newline, or quote, wrap in quotes and escape internal quotes
	if strings.Contains(field, ",") || strings.Contains(field, "\n") || strings.Contains(field, "\"") {
		escaped := strings.ReplaceAll(field, "\"", "\"\"")
		return "\"" + escaped + "\""
	}
	return field
}

func extractShellName(shellPath string) string {
	if shellPath == "" {
		return "unknown"
	}
	
	// Extract shell name from path (e.g., "/bin/bash" -> "bash")
	parts := strings.Split(shellPath, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return shellPath
}