package main
// package goutils

import (
	"bytes"
	"io"
	"os"
	"testing"
    "fmt"
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

func TestRichPrint(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Standard colors",
			input:    "[red on blue]Hello[/]",
			expected: "\x1b[31;44mHello\x1b[0m\n",
		},
		{
			name:     "Hex foreground",
			input:    "[#ff0000]Red[/]",
			expected: "\x1b[38;2;255;0;0mRed\x1b[0m\n",
		},
		{
			name:     "Hex background",
			input:    "[white on #00ff00]Green BG[/]",
			expected: "\x1b[37;48;2;0;255;0mGreen BG\x1b[0m\n",
		},
		{
			name:     "Bold attribute",
			input:    "[red bold]Bold Red[/]",
			expected: "\x1b[31;1mBold Red\x1b[0m\n",
		},
		{
			name:     "Hex foreground with standard background",
			input:    "[#ffffff on green]White on Green[/]",
			expected: "\x1b[38;2;255;255;255;42mWhite on Green\x1b[0m\n",
		},
		{
			name:     "Multiple styles",
			input:    "[red on blue bold]Hello[/] [green]World[/]",
			expected: "\x1b[31;44;1mHello\x1b[0m \x1b[32mWorld\x1b[0m\n",
		},
	}

    fmt.Printf("\nprint_test.go\n\n")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(func() {
				richPrint(tt.input, false) // Set debug to false for tests
			})
            message := "richPrint(%q) = `%q`; want `%q`"
			if output != tt.expected {
                t.Errorf(message, tt.input, output, tt.expected)
				return
			}
            fmt.Printf("    [PASS] " + message + "\n", tt.input, output, tt.expected)
		})
	}
    fmt.Printf("\n")
}
