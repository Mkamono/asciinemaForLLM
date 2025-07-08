package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/Mkamono/asciinemaForLLM/internal/formatter"
	"github.com/Mkamono/asciinemaForLLM/internal/parser"
)

// RunFormat reads from stdin and formats the output
func RunFormat(outputFormat string) error {
	scanner := bufio.NewScanner(os.Stdin)

	var lines []string
	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading input: %w", err)
	}

	session, err := parser.ParseCastFile(lines)
	if err != nil {
		return fmt.Errorf("error parsing cast file: %w", err)
	}

	if session == nil {
		return fmt.Errorf("no session data found")
	}

	switch outputFormat {
	case "csv":
		formatter.FormatSessionAsCSV(session)
	case "structured", "":
		formatter.FormatSessionAsStructured(session)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}
	
	return nil
}

// RunRecord starts asciinema recording, processes the output, and cleans up
func RunRecord(outputFile string, cleanup bool, outputFormat string) error {
	// Generate temporary filename if not provided
	if outputFile == "" {
		outputFile = fmt.Sprintf("session_%d.cast", time.Now().Unix())
	}

	// Ensure .cast extension
	if !strings.HasSuffix(outputFile, ".cast") {
		outputFile += ".cast"
	}

	fmt.Printf("Starting asciinema recording...\n")
	fmt.Printf("Recording will be saved to: %s\n", outputFile)
	fmt.Printf("Press Ctrl+D or type 'exit' to stop recording.\n\n")

	// Start asciinema recording
	cmd := exec.Command("asciinema", "rec", outputFile)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("asciinema recording failed: %w", err)
	}

	fmt.Printf("\nRecording completed: %s\n", outputFile)

	// Process the recorded file
	if err := processRecordedFile(outputFile, outputFormat); err != nil {
		return fmt.Errorf("failed to process recorded file: %w", err)
	}

	// Clean up original file if requested
	if cleanup {
		if err := os.Remove(outputFile); err != nil {
			fmt.Printf("Warning: failed to remove original file %s: %v\n", outputFile, err)
		} else {
			fmt.Printf("Original file %s removed.\n", outputFile)
		}
	}

	return nil
}

// RunFormatFile formats an existing .cast file
func RunFormatFile(inputFile string, outputFile string, cleanup bool, outputFormat string) error {
	// Check if input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// Generate output filename if not provided
	if outputFile == "" {
		base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
		if outputFormat == "csv" {
			outputFile = base + "_formatted.csv"
		} else {
			outputFile = base + "_formatted.md"
		}
	}

	// Process the file
	if err := processFile(inputFile, outputFile, outputFormat); err != nil {
		return fmt.Errorf("failed to process file: %w", err)
	}

	// Clean up original file if requested
	if cleanup {
		if err := os.Remove(inputFile); err != nil {
			fmt.Printf("Warning: failed to remove original file %s: %v\n", inputFile, err)
		} else {
			fmt.Printf("Original file %s removed.\n", inputFile)
		}
	}

	return nil
}

func processRecordedFile(inputFile string, outputFormat string) error {
	// Generate output filename
	base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
	var outputFile string
	if outputFormat == "csv" {
		outputFile = base + "_formatted.csv"
	} else {
		outputFile = base + "_formatted.md"
	}

	return processFile(inputFile, outputFile, outputFormat)
}

func processFile(inputFile, outputFile, outputFormat string) error {
	// Read the cast file
	file, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file: %w", err)
	}

	// Parse the cast file
	session, err := parser.ParseCastFile(lines)
	if err != nil {
		return fmt.Errorf("error parsing cast file: %w", err)
	}

	if session == nil {
		return fmt.Errorf("no session data found")
	}

	// Create output file
	outFile, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	// Redirect stdout to file temporarily
	oldStdout := os.Stdout
	os.Stdout = outFile

	// Format the session
	switch outputFormat {
	case "csv":
		formatter.FormatSessionAsCSV(session)
	case "structured", "":
		formatter.FormatSessionAsStructured(session)
	default:
		return fmt.Errorf("unknown output format: %s", outputFormat)
	}

	// Restore stdout
	os.Stdout = oldStdout

	fmt.Printf("Formatted output saved to: %s\n", outputFile)
	return nil
}

// ShowUsage displays help information
func ShowUsage() {
	fmt.Println(`asciinema-for-llm - asciinema録画ファイルをLLM向けに変換

使用方法:
    asciinema-for-llm [コマンド] [オプション]

コマンド:
    format              標準入力から.castファイルを読み取り、フォーマット済みテキストを出力
    record [ファイル名]   asciinema録画を開始し、終了後に自動でフォーマット
    file <入力> [出力]   既存の.castファイルをフォーマット

オプション:
    -h, --help          このヘルプメッセージを表示
    --cleanup           処理後に元の.castファイルを削除（record、fileコマンドで使用可能）
    --output=FORMAT     出力形式を指定（structured|csv、デフォルト: structured）

使用例:
    # 最も実用的な使い方（推奨）
    asciinema-for-llm record my_session.cast --output=csv --cleanup

    # 既存ファイルをCSV変換
    asciinema-for-llm file demo.cast --output=csv --cleanup

    # 標準入力からフォーマット
    cat demo.cast | asciinema-for-llm format --output=csv

    # 構造化テキスト形式（人間向け）
    asciinema-for-llm record my_session.cast --cleanup`)
}