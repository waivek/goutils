package goutils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	f()

	w.Close()
	os.Stdout = old
	return <-outC
}

func TestRichSprintf(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard colors",
			input:    "[red on blue]Hello[/]",
			expected: "\x1b[31;44mHello\x1b[0m",
		},
		{
			name:     "Hex foreground",
			input:    "[#ff0000]Red[/]",
			expected: "\x1b[38;2;255;0;0mRed\x1b[0m",
		},
		{
			name:     "Hex background",
			input:    "[white on #00ff00]Green BG[/]",
			expected: "\x1b[37;48;2;0;255;0mGreen BG\x1b[0m",
		},
		{
			name:     "Bold attribute",
			input:    "[red bold]Bold Red[/]",
			expected: "\x1b[31;1mBold Red\x1b[0m",
		},
		{
			name:     "Hex foreground with standard background",
			input:    "[#ffffff on green]White on Green[/]",
			expected: "\x1b[38;2;255;255;255;42mWhite on Green\x1b[0m",
		},
		{
			name:     "Multiple styles",
			input:    "[red on blue bold]Hello[/] [green]World[/]",
			expected: "\x1b[31;44;1mHello\x1b[0m \x1b[32mWorld\x1b[0m",
		},
	}

	fmt.Printf("\nprint_test.go\n\n")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := richSprintf(tt.input)
            message := fmt.Sprintf("richSprintf(%q) = %q, want %q", tt.input, output, tt.expected)
			if output != tt.expected {
				t.Errorf(message)
			}
            pass_string_colored := richSprintf("[black on green bold] PASS [/]")
            // fmt.Printf("    PASS: %s\n", message)
            fmt.Printf("    %s %s\n", pass_string_colored, message)
		})
	}
	fmt.Printf("\n")
}

