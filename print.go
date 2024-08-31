package goutils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

func richPrint(text string, debugOpt ...bool) {
	debug := false
	if len(debugOpt) > 0 {
		debug = debugOpt[0]
	}

	formattedText := richSprintf(text)

	if debug {
		fmt.Fprintf(os.Stderr, "Original text: %s\n", text)
		visibleAnsi := strings.ReplaceAll(formattedText, "\x1b", "\\x1b")
		fmt.Fprintf(os.Stderr, "Final text with visible ANSI: %s\n", visibleAnsi)
	}

	fmt.Println(formattedText)
}

func parseAttributes(attrs []string) []color.Attribute {
	result := []color.Attribute{}
	for _, attr := range attrs {
		switch attr {
		case "bold":
			result = append(result, color.Bold)
			// Add more attributes as needed
		}
	}
	return result
}

func createColorPrinter(fg, bg string, attrs []color.Attribute) func(a ...interface{}) string {
	fgColor, fgIsHex := parseColor(fg)
	var bgColor color.Attribute
	var bgIsHex bool

	if bg != "" {
		bgColor, bgIsHex = parseColor(bg)
		if !bgIsHex {
			bgColor += 10 // Convert standard color to background
		}
	}

	// Convert attributes to ANSI codes
	attrCodes := ""
	for _, attr := range attrs {
		attrCodes += fmt.Sprintf(";%d", attr)
	}

	if bg == "" {
		if fgIsHex {
			r, g, b := hexToRGB(fg)
			return func(a ...interface{}) string {
				return fmt.Sprintf("\x1b[38;2;%d;%d;%d%sm%s\x1b[0m", r, g, b, attrCodes, fmt.Sprint(a...))
			}
		}
		return func(a ...interface{}) string {
			return fmt.Sprintf("\x1b[%d%sm%s\x1b[0m", fgColor, attrCodes, fmt.Sprint(a...))
		}
	}

	if fgIsHex {
		fgR, fgG, fgB := hexToRGB(fg)
		if bgIsHex {
			bgR, bgG, bgB := hexToRGB(bg)
			return func(a ...interface{}) string {
				return fmt.Sprintf("\x1b[38;2;%d;%d;%d;48;2;%d;%d;%d%sm%s\x1b[0m", fgR, fgG, fgB, bgR, bgG, bgB, attrCodes, fmt.Sprint(a...))
			}
		}
		return func(a ...interface{}) string {
			return fmt.Sprintf("\x1b[38;2;%d;%d;%d;%d%sm%s\x1b[0m", fgR, fgG, fgB, bgColor, attrCodes, fmt.Sprint(a...))
		}
	}

	if bgIsHex {
		bgR, bgG, bgB := hexToRGB(bg)
		return func(a ...interface{}) string {
			return fmt.Sprintf("\x1b[%d;48;2;%d;%d;%d%sm%s\x1b[0m", fgColor, bgR, bgG, bgB, attrCodes, fmt.Sprint(a...))
		}
	}

	return func(a ...interface{}) string {
		return fmt.Sprintf("\x1b[%d;%d%sm%s\x1b[0m", fgColor, bgColor, attrCodes, fmt.Sprint(a...))
	}
}

func parseColor(colorName string) (color.Attribute, bool) {
	switch strings.ToLower(colorName) {
	case "black":
		return color.FgBlack, false
	case "red":
		return color.FgRed, false
	case "green":
		return color.FgGreen, false
	case "yellow":
		return color.FgYellow, false
	case "blue":
		return color.FgBlue, false
	case "magenta":
		return color.FgMagenta, false
	case "cyan":
		return color.FgCyan, false
	case "white":
		return color.FgWhite, false
	default:
		if strings.HasPrefix(colorName, "#") {
			return color.Attribute(38), true // 38 for foreground, true for isHex
		}
		return color.FgWhite, false
	}
}

func hexToRGB(hex string) (uint8, uint8, uint8) {
	hex = strings.TrimPrefix(hex, "#")
	if len(hex) == 3 {
		hex = string([]byte{hex[0], hex[0], hex[1], hex[1], hex[2], hex[2]})
	}
	r, _ := strconv.ParseUint(hex[0:2], 16, 8)
	g, _ := strconv.ParseUint(hex[2:4], 16, 8)
	b, _ := strconv.ParseUint(hex[4:6], 16, 8)
	return uint8(r), uint8(g), uint8(b)
}

func richSprintf(text string) string {
	re := regexp.MustCompile(`\[([^\]]+)\]([^\[]+)\[/\]`)
	matches := re.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		fullMatch := match[0]
		style := match[1]
		content := match[2]

		parts := strings.Split(style, " ")
		fg := parts[0]
		bg := ""
		attrs := []color.Attribute{}

		if len(parts) > 2 && parts[1] == "on" {
			bg = parts[2]
			attrs = parseAttributes(parts[3:])
		} else {
			attrs = parseAttributes(parts[1:])
		}

		printer := createColorPrinter(fg, bg, attrs)
		coloredText := printer(content)

		text = strings.Replace(text, fullMatch, coloredText, 1)
	}

	return text
}

func main() {
	richPrint("[red on blue bold]hello[/] [red on #ffffff]world[/]")
	richPrint("[black on green bold]hello[/]")
	richPrint("[white on green bold]hello[/]")
	richPrint("[black on green]hello[/]")
	richPrint("[white on green]hello[/]")
	richPrint("[#ffffff on green] HELLO [/]")
	richPrint("[#ffffff on green bold] HELLO [/]")
}
