package main

import (
	"bufio"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/thalessoares/lg/internal/buffer"
	"github.com/thalessoares/lg/internal/parser"
	"github.com/thalessoares/lg/internal/tui"
)

const (
	bufferCapacity = 10000
)

func main() {
	// Check if stdin is a pipe
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		fmt.Fprintln(os.Stderr, "Usage: <command> | lg")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "lg reads JSON logs from stdin and displays them in an interactive TUI.")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Example:")
		fmt.Fprintln(os.Stderr, "  tail -f app.log | lg")
		fmt.Fprintln(os.Stderr, "  docker logs -f container | lg")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Keybindings:")
		fmt.Fprintln(os.Stderr, "  j/k, arrows  : scroll up/down")
		fmt.Fprintln(os.Stderr, "  g/G          : go to top/bottom")
		fmt.Fprintln(os.Stderr, "  Ctrl+d/u     : page down/up")
		fmt.Fprintln(os.Stderr, "  /            : search/filter")
		fmt.Fprintln(os.Stderr, "  p            : pause/resume")
		fmt.Fprintln(os.Stderr, "  c            : clear logs")
		fmt.Fprintln(os.Stderr, "  q, Ctrl+c    : quit")
		os.Exit(1)
	}

	// Create buffer
	buf := buffer.New(bufferCapacity)

	// Create TUI model
	model := tui.New(buf)

	// Create program with stdin reading disabled (we'll read from stdin ourselves)
	p := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Start reading stdin in a goroutine
	go func() {
		scanner := bufio.NewScanner(os.Stdin)
		// Increase buffer size for long lines
		const maxScanTokenSize = 1024 * 1024 // 1MB
		scanBuf := make([]byte, maxScanTokenSize)
		scanner.Buffer(scanBuf, maxScanTokenSize)

		for scanner.Scan() {
			line := scanner.Text()
			if entry := parser.Parse(line); entry != nil {
				p.Send(tui.AddLogEntry(entry))
			}
		}

		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		}
	}()

	// Run the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}
