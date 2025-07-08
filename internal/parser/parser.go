package parser

import (
	"encoding/json"
	"regexp"
	"strconv"
	"strings"
)

type CastHeader struct {
	Version   int               `json:"version"`
	Width     int               `json:"width"`
	Height    int               `json:"height"`
	Timestamp int64             `json:"timestamp"`
	Env       map[string]string `json:"env"`
}

type CastEvent struct {
	Timestamp float64 `json:"timestamp"`
	EventType string  `json:"event_type"`
	Data      string  `json:"data"`
}

type Command struct {
	Command    string
	StartTime  float64
	EndTime    float64
	Output     string
	IsComplete bool
	ExitCode   int
}

type CastSession struct {
	Header   CastHeader
	Commands []Command
}

func ParseCastFile(lines []string) (*CastSession, error) {
	if len(lines) == 0 {
		return nil, nil
	}

	var header CastHeader
	var events []CastEvent

	// Parse header
	if err := json.Unmarshal([]byte(lines[0]), &header); err != nil {
		return nil, err
	}

	// Parse events
	for i := 1; i < len(lines); i++ {
		line := lines[i]
		if line == "" {
			continue
		}

		var event []interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}

		if len(event) >= 3 {
			timestamp, _ := event[0].(float64)
			eventType, _ := event[1].(string)
			data, _ := event[2].(string)

			events = append(events, CastEvent{
				Timestamp: timestamp,
				EventType: eventType,
				Data:      data,
			})
		}
	}

	commands := parseCommands(events)

	return &CastSession{
		Header:   header,
		Commands: commands,
	}, nil
}

func parseCommands(events []CastEvent) []Command {
	var commands []Command
	var currentCommand Command
	var buffer strings.Builder
	commandExecuted := false
	outputStarted := false

	for _, event := range events {
		if event.EventType != "o" {
			continue
		}

		cleaned := cleanTerminalOutput(event.Data)

		// Check for exit status
		if exitCode := extractExitCode(event.Data); exitCode != -1 && currentCommand.Command != "" {
			currentCommand.ExitCode = exitCode
		}

		if isCommandExecution(event.Data) {
			if currentCommand.Command != "" {
				currentCommand.EndTime = event.Timestamp
				output := strings.TrimSpace(buffer.String())
				currentCommand.Output = cleanFinalOutput(output)
				currentCommand.IsComplete = true
				commands = append(commands, currentCommand)
				buffer.Reset()
			}

			cmd := extractCommandFromExecution(event.Data)
			if cmd != "" {
				currentCommand = Command{
					Command:    cmd,
					StartTime:  event.Timestamp,
					IsComplete: false,
					ExitCode:   0, // Default to success, will be updated if status found
				}
				commandExecuted = true
				outputStarted = false
			}
		} else if commandExecuted && currentCommand.Command != "" {
			if isActualOutput(cleaned) {
				outputStarted = true
			}

			if cleaned != "" && !isOnlyControlSequence(cleaned) && !isPromptOutput(cleaned) {
				if outputStarted && !isJunkOutput(cleaned) {
					if buffer.Len() > 0 {
						buffer.WriteString("\n")
					}
					buffer.WriteString(cleaned)
				} else if isActualOutput(cleaned) {
					outputStarted = true
					if buffer.Len() > 0 {
						buffer.WriteString("\n")
					}
					buffer.WriteString(cleaned)
				}
			}
			currentCommand.EndTime = event.Timestamp
		}
	}

	if currentCommand.Command != "" {
		output := strings.TrimSpace(buffer.String())
		currentCommand.Output = cleanFinalOutput(output)
		currentCommand.IsComplete = true
		commands = append(commands, currentCommand)
	}

	return commands
}

func cleanTerminalOutput(data string) string {
	// Remove ANSI escape sequences
	ansiRegex := regexp.MustCompile(`\x1b\[[0-9;]*[a-zA-Z]`)
	cleaned := ansiRegex.ReplaceAllString(data, "")

	// Remove OSC sequences (e.g., \x1b]0;title\x07)
	oscRegex := regexp.MustCompile(`\x1b\][^\\x07]*\x07`)
	cleaned = oscRegex.ReplaceAllString(cleaned, "")

	// Remove other control sequences
	controlRegex := regexp.MustCompile(`\x1b[\[\]><=?][^a-zA-Z]*[a-zA-Z]`)
	cleaned = controlRegex.ReplaceAllString(cleaned, "")

	// Remove specific escape sequences
	cleaned = strings.ReplaceAll(cleaned, "\x1b(B", "")
	cleaned = strings.ReplaceAll(cleaned, "\x1b[m", "")
	cleaned = strings.ReplaceAll(cleaned, "\x1b=", "")
	cleaned = strings.ReplaceAll(cleaned, "\x1b>", "")

	// Remove control characters
	cleaned = strings.ReplaceAll(cleaned, "\x00", "")
	cleaned = strings.ReplaceAll(cleaned, "\x07", "")
	cleaned = strings.ReplaceAll(cleaned, "\x08", "")
	cleaned = strings.ReplaceAll(cleaned, "\x0c", "")
	cleaned = strings.ReplaceAll(cleaned, "\x0e", "")
	cleaned = strings.ReplaceAll(cleaned, "\x0f", "")

	// Normalize line endings
	cleaned = regexp.MustCompile(`\r\n|\r`).ReplaceAllString(cleaned, "\n")

	return strings.TrimSpace(cleaned)
}

func isCommandExecution(data string) bool {
	return strings.Contains(data, "cmdline_url=")
}

func extractCommandFromExecution(data string) string {
	cmdRegex := regexp.MustCompile(`cmdline_url=([^\x07]+)`)
	matches := cmdRegex.FindStringSubmatch(data)
	if len(matches) > 1 {
		decoded := strings.ReplaceAll(matches[1], "%20", " ")
		decoded = strings.ReplaceAll(decoded, "%22", "\"")
		decoded = strings.ReplaceAll(decoded, "%2C", ",")
		return decoded
	}
	return ""
}

func isOnlyControlSequence(text string) bool {
	if text == "" {
		return true
	}
	cleanedForCheck := regexp.MustCompile(`[^\x00-\x1f\x7f-\x9f]`).ReplaceAllString(text, "")
	return len(cleanedForCheck) == len(text)
}

func isActualOutput(text string) bool {
	// Check for actual command output patterns - not prompt or input
	// Look for meaningful text that could be command output
	return text != "" &&
		!strings.Contains(text, "❯") &&
		!strings.Contains(text, "$") &&
		!strings.Contains(text, "#") &&
		!isPromptOutput(text) &&
		(len(text) > 5 || // longer text is likely output
			regexp.MustCompile(`^[A-Z]`).MatchString(text) || // starts with capital
			strings.Contains(text, " ") || // contains spaces
			strings.HasPrefix(text, "/")) // file paths are valid output
}

func isPromptOutput(text string) bool {
	// Check if text contains prompt indicators
	return strings.Contains(text, "❯") ||
		strings.Contains(text, "kamonomakoto@") ||
		strings.Contains(text, "~/") ||
		strings.Contains(text, "(main)")
}

func isJunkOutput(text string) bool {
	// Filter out junk output like random characters, typing artifacts, etc.
	return strings.Contains(text, "cho \"") ||
		strings.Contains(text, "echo \"") ||
		strings.Contains(text, "ec fish") ||
		strings.Contains(text, "⏎") ||
		strings.Contains(text, "/r/asciinemaForLLM") ||
		(regexp.MustCompile(`^[a-z]$`).MatchString(text)) ||
		(regexp.MustCompile(`^[a-z]{1,4}$`).MatchString(text)) || // includes pwd, exit, pipx
		isCommandTyping(text)
}

func isCommandTyping(text string) bool {
	// Detect command typing artifacts that appear during command input
	commonCommands := []string{"ls", "pwd", "exit", "echo", "cat", "cd", "mkdir", "rm", "cp", "mv"}
	for _, cmd := range commonCommands {
		if text == cmd {
			return true
		}
	}
	return false
}

func cleanFinalOutput(output string) string {
	if output == "" {
		return ""
	}

	// Split by lines and filter out unwanted lines
	lines := strings.Split(output, "\n")
	var cleanLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !isPromptOutput(line) && !isJunkOutput(line) && !strings.Contains(line, ";0") {
			// Include any legitimate output that's not prompt or junk
			cleanLines = append(cleanLines, line)
		}
	}

	return strings.Join(cleanLines, "\n")
}

func extractExitCode(data string) int {
	// Look for fish shell exit status format: \u001b]133;D;code\u0007
	statusRegex := regexp.MustCompile(`\x1b\]133;D;(\d+)\x07`)
	matches := statusRegex.FindStringSubmatch(data)
	if len(matches) > 1 {
		if code, err := strconv.Atoi(matches[1]); err == nil {
			return code
		}
	}
	
	// Also look for alternative status_url format (if any shell uses it)
	statusURLRegex := regexp.MustCompile(`status_url=(\d+)`)
	matches = statusURLRegex.FindStringSubmatch(data)
	if len(matches) > 1 {
		if code, err := strconv.Atoi(matches[1]); err == nil {
			return code
		}
	}
	
	return -1 // No exit code found
}