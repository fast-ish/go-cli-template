// Package output provides formatted output for the CLI
package output

import (
	"encoding/json"
	"fmt"
	"os"
{%- if values.outputFormat == "charm" %}
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/huh"
{%- elif values.outputFormat == "tablewriter" %}
	"github.com/olekukonko/tablewriter"
{%- endif %}
{%- if values.configFormat == "yaml" or values.configFormat == "all" %}
	"gopkg.in/yaml.v3"
{%- endif %}
)

{%- if values.outputFormat == "charm" %}

// Styles using lipgloss
var (
	StyleTitle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	StyleSuccess = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	StyleError   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	StyleWarning = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	StyleInfo    = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
)
{%- endif %}

// Formatter handles output formatting
type Formatter struct {
	format string
	color  bool
}

// NewFormatter creates a new output formatter
func NewFormatter(format string) *Formatter {
	return &Formatter{
		format: format,
		color:  true, // TODO: read from global flags
	}
}

// Data outputs data in the configured format
func (f *Formatter) Data(data any, title string) error {
	switch f.format {
	case "json":
		return f.JSON(data)
	case "yaml":
		return f.YAML(data)
	case "table":
		return f.Table(data, title)
	default:
		return f.Auto(data, title)
	}
}

// JSON outputs data as JSON
func (f *Formatter) JSON(data any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// YAML outputs data as YAML
func (f *Formatter) YAML(data any) error {
{%- if values.configFormat == "yaml" or values.configFormat == "all" %}
	enc := yaml.NewEncoder(os.Stdout)
	defer enc.Close()
	return enc.Encode(data)
{%- else %}
	return fmt.Errorf("YAML output not supported (enable in config)")
{%- endif %}
}

// Table outputs data as a table
func (f *Formatter) Table(data any, title string) error {
	// Convert to []map[string]any for table rendering
	items, ok := data.([]map[string]any)
	if !ok {
		return fmt.Errorf("table format requires []map[string]any")
	}

	if len(items) == 0 {
		f.Info("No items to display")
		return nil
	}

{%- if values.outputFormat == "charm" %}
	// Extract headers from first item
	var headers []string
	for k := range items[0] {
		headers = append(headers, k)
	}

	// Build rows
	var rows []table.Row
	for _, item := range items {
		var row table.Row
		for _, h := range headers {
			row = append(row, fmt.Sprintf("%v", item[h]))
		}
		rows = append(rows, row)
	}

	// Create table columns
	var columns []table.Column
	for _, h := range headers {
		columns = append(columns, table.Column{Title: h, Width: 20})
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(rows)),
	)

	if title != "" {
		fmt.Println(StyleTitle.Render(title))
	}
	fmt.Println(t.View())
{%- elif values.outputFormat == "tablewriter" %}
	// Extract headers
	var headers []string
	for k := range items[0] {
		headers = append(headers, k)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)

	for _, item := range items {
		var row []string
		for _, h := range headers {
			row = append(row, fmt.Sprintf("%v", item[h]))
		}
		table.Append(row)
	}

	if title != "" {
		fmt.Println(title)
	}
	table.Render()
{%- else %}
	// Plain text table
	if title != "" {
		fmt.Println(title)
	}
	for _, item := range items {
		for k, v := range item {
			fmt.Printf("%s: %v\n", k, v)
		}
		fmt.Println()
	}
{%- endif %}

	return nil
}

// Auto automatically selects format based on data type
func (f *Formatter) Auto(data any, title string) error {
	switch v := data.(type) {
	case []map[string]any:
		return f.Table(v, title)
	default:
		return f.JSON(data)
	}
}

// Success prints a success message
func (f *Formatter) Success(msg string) {
{%- if values.outputFormat == "charm" %}
	fmt.Println(StyleSuccess.Render("✓ " + msg))
{%- else %}
	fmt.Println("✓", msg)
{%- endif %}
}

// Error prints an error message
func (f *Formatter) Error(msg string) {
{%- if values.outputFormat == "charm" %}
	fmt.Fprintln(os.Stderr, StyleError.Render("✗ "+msg))
{%- else %}
	fmt.Fprintln(os.Stderr, "✗", msg)
{%- endif %}
}

// Warning prints a warning message
func (f *Formatter) Warning(msg string) {
{%- if values.outputFormat == "charm" %}
	fmt.Println(StyleWarning.Render("⚠ " + msg))
{%- else %}
	fmt.Println("⚠", msg)
{%- endif %}
}

// Info prints an info message
func (f *Formatter) Info(msg string) {
{%- if values.outputFormat == "charm" %}
	fmt.Println(StyleInfo.Render("ℹ " + msg))
{%- else %}
	fmt.Println("ℹ", msg)
{%- endif %}
}

// Confirm prompts the user for confirmation
func (f *Formatter) Confirm(message string, defaultValue bool) bool {
{%- if values.outputFormat == "charm" %}
	var confirm bool
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(message).
				Value(&confirm).
				Affirmative("Yes").
				Negative("No"),
		),
	)

	err := form.Run()
	if err != nil {
		return defaultValue
	}
	return confirm
{%- else %}
	// Simple text-based confirmation
	var response string
	defaultStr := "n"
	if defaultValue {
		defaultStr = "y"
	}
	fmt.Printf("%s [y/N] (default: %s): ", message, defaultStr)
	fmt.Scanln(&response)

	if response == "" {
		return defaultValue
	}
	return response == "y" || response == "Y" || response == "yes"
{%- endif %}
}

// DryRun prints what would happen in dry-run mode
func (f *Formatter) DryRun(format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
{%- if values.outputFormat == "charm" %}
	fmt.Println(StyleWarning.Render("[DRY RUN] " + msg))
{%- else %}
	fmt.Println("[DRY RUN]", msg)
{%- endif %}
}
