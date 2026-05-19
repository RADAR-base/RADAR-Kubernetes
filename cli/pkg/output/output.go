package output

import (
	"encoding/json"
	"os"

	"github.com/pterm/pterm"
)

type Mode string

const (
	Human Mode = "human"
	JSON  Mode = "json"
)

// PrintJSON writes v as indented JSON to stdout.
func PrintJSON(v any) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(v)
}

// Success prints a green checkmark line.
func Success(msg string) {
	pterm.Success.Println(msg)
}

// Error prints a red error line.
func Error(msg string) {
	pterm.Error.Println(msg)
}

// Warning prints a yellow warning line.
func Warning(msg string) {
	pterm.Warning.Println(msg)
}

// Info prints an info line.
func Info(msg string) {
	pterm.Info.Println(msg)
}

// Header prints a bold section header.
func Header(msg string) {
	pterm.DefaultHeader.WithFullWidth().Println(msg)
}

// ValidationJSON is the JSON schema for validate output.
type ValidationJSON struct {
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

