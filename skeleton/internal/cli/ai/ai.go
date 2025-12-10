{%- if values.aiProvider != "none" %}
// Package ai provides AI command-line interface
package ai

import (
{%- if values.cliFramework == "cobra" %}
	"github.com/spf13/cobra"
{%- elif values.cliFramework == "urfave" %}
	"github.com/urfave/cli/v2"
{%- endif %}

	"github.com/fast-ish/${{values.name}}/internal/context"
)

{%- if values.cliFramework == "cobra" %}

// Cmd is the root AI command
var Cmd = &cobra.Command{
	Use:   "ai",
	Short: "AI-powered operations",
	Long:  "AI-powered operations using {{values.aiProvider}}",
}

func init() {
{%- if "chat" in values.aiFeatures %}
	Cmd.AddCommand(chatCmd)
{%- endif %}
{%- if "analyze" in values.aiFeatures %}
	Cmd.AddCommand(analyzeCmd)
{%- endif %}
{%- if "summarize" in values.aiFeatures %}
	Cmd.AddCommand(summarizeCmd)
{%- endif %}
{%- if "generate" in values.aiFeatures %}
	Cmd.AddCommand(generateCmd)
{%- endif %}
	Cmd.AddCommand(modelsCmd)
}

{%- if "chat" in values.aiFeatures %}

var chatCmd = &cobra.Command{
	Use:   "chat [prompt]",
	Short: "Chat with AI",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.GetGlobal()
		prompt := args[0]

		response, err := ctx.AI().Chat(cmd.Context(), prompt)
		if err != nil {
			return err
		}

		ctx.Output.Info(response)
		return nil
	},
}
{%- endif %}

{%- if "analyze" in values.aiFeatures %}

var analyzeCmd = &cobra.Command{
	Use:   "analyze [text]",
	Short: "Analyze text with AI",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.GetGlobal()
		text := args[0]

		analysis, err := ctx.AI().Analyze(cmd.Context(), text)
		if err != nil {
			return err
		}

		ctx.Output.Info(analysis)
		return nil
	},
}
{%- endif %}

{%- if "summarize" in values.aiFeatures %}

var summarizeCmd = &cobra.Command{
	Use:   "summarize [text]",
	Short: "Summarize text with AI",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.GetGlobal()
		text := args[0]

		summary, err := ctx.AI().Summarize(cmd.Context(), text)
		if err != nil {
			return err
		}

		ctx.Output.Info(summary)
		return nil
	},
}
{%- endif %}

{%- if "generate" in values.aiFeatures %}

var generateCmd = &cobra.Command{
	Use:   "generate [prompt]",
	Short: "Generate content with AI",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.GetGlobal()
		prompt := args[0]

		content, err := ctx.AI().Generate(cmd.Context(), prompt, nil)
		if err != nil {
			return err
		}

		ctx.Output.Info(content)
		return nil
	},
}
{%- endif %}

var modelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List available AI models",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.GetGlobal()
		ctx.Output.Info("Current model: " + ctx.Config.AI.Model)
		return nil
	},
}

{%- elif values.cliFramework == "urfave" %}

// Cmd is the root AI command
var Cmd = &cli.Command{
	Name:  "ai",
	Usage: "AI-powered operations using {{values.aiProvider}}",
	Subcommands: []*cli.Command{
{%- if "chat" in values.aiFeatures %}
		chatCmd,
{%- endif %}
{%- if "analyze" in values.aiFeatures %}
		analyzeCmd,
{%- endif %}
{%- if "summarize" in values.aiFeatures %}
		summarizeCmd,
{%- endif %}
{%- if "generate" in values.aiFeatures %}
		generateCmd,
{%- endif %}
		modelsCmd,
	},
}

{%- if "chat" in values.aiFeatures %}

var chatCmd = &cli.Command{
	Name:      "chat",
	Usage:     "Chat with AI",
	ArgsUsage: "[prompt]",
	Action: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.ShowSubcommandHelp(c)
		}

		ctx := context.GetGlobal()
		prompt := c.Args().First()

		response, err := ctx.AI().Chat(c.Context, prompt)
		if err != nil {
			return err
		}

		ctx.Output.Info(response)
		return nil
	},
}
{%- endif %}

{%- if "analyze" in values.aiFeatures %}

var analyzeCmd = &cli.Command{
	Name:      "analyze",
	Usage:     "Analyze text with AI",
	ArgsUsage: "[text]",
	Action: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.ShowSubcommandHelp(c)
		}

		ctx := context.GetGlobal()
		text := c.Args().First()

		analysis, err := ctx.AI().Analyze(c.Context, text)
		if err != nil {
			return err
		}

		ctx.Output.Info(analysis)
		return nil
	},
}
{%- endif %}

{%- if "summarize" in values.aiFeatures %}

var summarizeCmd = &cli.Command{
	Name:      "summarize",
	Usage:     "Summarize text with AI",
	ArgsUsage: "[text]",
	Action: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.ShowSubcommandHelp(c)
		}

		ctx := context.GetGlobal()
		text := c.Args().First()

		summary, err := ctx.AI().Summarize(c.Context, text)
		if err != nil {
			return err
		}

		ctx.Output.Info(summary)
		return nil
	},
}
{%- endif %}

{%- if "generate" in values.aiFeatures %}

var generateCmd = &cli.Command{
	Name:      "generate",
	Usage:     "Generate content with AI",
	ArgsUsage: "[prompt]",
	Action: func(c *cli.Context) error {
		if c.NArg() < 1 {
			return cli.ShowSubcommandHelp(c)
		}

		ctx := context.GetGlobal()
		prompt := c.Args().First()

		content, err := ctx.AI().Generate(c.Context, prompt, nil)
		if err != nil {
			return err
		}

		ctx.Output.Info(content)
		return nil
	},
}
{%- endif %}

var modelsCmd = &cli.Command{
	Name:  "models",
	Usage: "List available AI models",
	Action: func(c *cli.Context) error {
		ctx := context.GetGlobal()
		ctx.Output.Info("Current model: " + ctx.Config.AI.Model)
		return nil
	},
}
{%- endif %}
{%- else %}
// Package ai is a placeholder when AI is not enabled
package ai

import (
{%- if values.cliFramework == "cobra" %}
	"github.com/spf13/cobra"

var Cmd = &cobra.Command{
	Use:   "ai",
	Short: "AI operations (disabled)",
	RunE: func(cmd *cobra.Command, args []string) error {
		return fmt.Errorf("AI features not enabled in this build")
	},
}
{%- elif values.cliFramework == "urfave" %}
	"github.com/urfave/cli/v2"
	"fmt"
)

var Cmd = &cli.Command{
	Name:  "ai",
	Usage: "AI operations (disabled)",
	Action: func(c *cli.Context) error {
		return fmt.Errorf("AI features not enabled in this build")
	},
}
{%- endif %}
{%- endif %}
