// Package cli provides the command-line interface for ${{values.name}}
package cli

import (
	"fmt"
	"os"

{%- if values.cliFramework == "cobra" %}
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
{%- elif values.cliFramework == "urfave" %}
	"github.com/urfave/cli/v2"
{%- endif %}

	"github.com/fast-ish/${{values.name}}/internal/config"
	"github.com/fast-ish/${{values.name}}/internal/context"
	"github.com/fast-ish/${{values.name}}/internal/logger"
	"github.com/fast-ish/${{values.name}}/internal/output"
{%- if values.aiProvider != "none" %}
	"github.com/fast-ish/${{values.name}}/internal/cli/ai"
{%- endif %}
	// Import command modules here
	// These are dynamically registered based on template selections
{%- for integration in values.integrations %}
	"github.com/fast-ish/${{values.name}}/internal/cli/{{integration}}"
{%- endfor %}
)

var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
)

// SetVersion sets the version information
func SetVersion(v, bt, gc string) {
	version = v
	buildTime = bt
	gitCommit = gc
}

{%- if values.cliFramework == "cobra" %}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "${{values.name}}",
	Short: "${{values.description}}",
	Long: `${{values.name}} - ${{values.description}}

A modular, extensible CLI tool built with the fast-ish golden path.

Features:
{%- if values.aiProvider != "none" %}
  • AI-powered operations with {{values.aiProvider}}
{%- endif %}
{%- if values.integrations|length > 0 %}
  • Integrations: {{values.integrations|join(", ")}}
{%- endif %}
  • Structured logging
  • Rich terminal output
{%- if values.metrics %}
  • Prometheus metrics
{%- endif %}
{%- if values.tracing %}
  • OpenTelemetry tracing
{%- endif %}
`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Initialize configuration
		cfg, err := config.Load(viper.GetString("config"))
		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Initialize global context
		ctx := context.NewContext(cfg)

		// Set verbosity
		verbose, _ := cmd.Flags().GetCount("verbose")
		ctx.Verbose = verbose

		// Set output format
		outputFormat, _ := cmd.Flags().GetString("output")
		ctx.Output = output.NewFormatter(outputFormat)

		// Set dry-run mode
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		ctx.DryRun = dryRun

		// Store in global context
		context.SetGlobal(ctx)

		return nil
	},
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default: ~/.{{values.name}}/config.yaml)")
	rootCmd.PersistentFlags().StringP("output", "o", "auto", "output format: auto, json, yaml, table")
	rootCmd.PersistentFlags().CountP("verbose", "v", "verbose output (-v for info, -vv for debug)")
	rootCmd.PersistentFlags().Bool("dry-run", false, "show what would happen without making changes")
	rootCmd.PersistentFlags().Bool("no-color", false, "disable colored output")

	// Bind flags to viper
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("output", rootCmd.PersistentFlags().Lookup("output"))
	viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))
	viper.BindPFlag("no-color", rootCmd.PersistentFlags().Lookup("no-color"))

	// Version command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version %s (built %s, commit %s)\n",
				"${{values.name}}", version, buildTime, gitCommit)
		},
	})

	// Config command
	rootCmd.AddCommand(&cobra.Command{
		Use:   "config",
		Short: "Show current configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.GetGlobal()
			return ctx.Output.Data(ctx.Config, "Configuration")
		},
	})

	// Register command modules
	registerCommands()
}

// registerCommands registers all command modules
// This is where the modular architecture shines - commands are auto-registered
func registerCommands() {
{%- if values.aiProvider != "none" %}
	rootCmd.AddCommand(ai.Cmd)
{%- endif %}
{%- for integration in values.integrations %}
	rootCmd.AddCommand({{integration}}.Cmd)
{%- endfor %}
}

{%- elif values.cliFramework == "urfave" %}

// Execute runs the CLI application
func Execute() error {
	app := &cli.App{
		Name:     "${{values.name}}",
		Usage:    "${{values.description}}",
		Version:  fmt.Sprintf("%s (built %s, commit %s)", version, buildTime, gitCommit),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "config file (default: ~/.{{values.name}}/config.yaml)",
			},
			&cli.StringFlag{
				Name:    "output",
				Aliases: []string{"o"},
				Value:   "auto",
				Usage:   "output format: auto, json, yaml, table",
			},
			&cli.IntFlag{
				Name:    "verbose",
				Aliases: []string{"v"},
				Value:   0,
				Usage:   "verbose output (0=warn, 1=info, 2=debug)",
			},
			&cli.BoolFlag{
				Name:  "dry-run",
				Usage: "show what would happen without making changes",
			},
			&cli.BoolFlag{
				Name:  "no-color",
				Usage: "disable colored output",
			},
		},
		Before: func(c *cli.Context) error {
			// Initialize configuration
			cfg, err := config.Load(c.String("config"))
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			// Initialize global context
			ctx := context.NewContext(cfg)
			ctx.Verbose = c.Int("verbose")
			ctx.Output = output.NewFormatter(c.String("output"))
			ctx.DryRun = c.Bool("dry-run")

			// Store in global context
			context.SetGlobal(ctx)

			return nil
		},
		Commands: registerCommands(),
	}

	return app.Run(os.Args)
}

// registerCommands registers all command modules
func registerCommands() []*cli.Command {
	commands := []*cli.Command{
		{
			Name:  "config",
			Usage: "Show current configuration",
			Action: func(c *cli.Context) error {
				ctx := context.GetGlobal()
				return ctx.Output.Data(ctx.Config, "Configuration")
			},
		},
	}

{%- if values.aiProvider != "none" %}
	commands = append(commands, ai.Cmd)
{%- endif %}
{%- for integration in values.integrations %}
	commands = append(commands, {{integration}}.Cmd)
{%- endfor %}

	return commands
}
{%- endif %}
