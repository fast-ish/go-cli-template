# Extending ${{values.name}}

This guide shows how to extend and customize the CLI with new features and modules.

## Table of Contents

- [Adding New Commands](#adding-new-commands)
- [Adding New Integrations](#adding-new-integrations)
- [Adding AI Features](#adding-ai-features)
- [Custom Output Formats](#custom-output-formats)
- [Plugin Development](#plugin-development)
- [Middleware and Hooks](#middleware-and-hooks)

---

## Adding New Commands

### 1. Create Command Package

Create a new package in `internal/cli/`:

```bash
mkdir -p internal/cli/myfeature
```

### 2. Define Command Structure

{%- if values.cliFramework == "cobra" %}

```go
// internal/cli/myfeature/myfeature.go
package myfeature

import (
    "fmt"

    "github.com/spf13/cobra"
    "github.com/fast-ish/${{values.name}}/internal/context"
    "github.com/fast-ish/${{values.name}}/internal/logger"
)

// Cmd is the root command for myfeature
var Cmd = &cobra.Command{
    Use:   "myfeature",
    Short: "My new feature commands",
    Long:  "Detailed description of what this feature does",
}

var doSomethingCmd = &cobra.Command{
    Use:   "do-something",
    Short: "Do something useful",
    Args:  cobra.ExactArgs(1),
    RunE:  runDoSomething,
}

func init() {
    // Register subcommands
    Cmd.AddCommand(doSomethingCmd)

    // Add flags
    doSomethingCmd.Flags().StringP("option", "o", "default", "An option")
    doSomethingCmd.Flags().BoolP("force", "f", false, "Force operation")
}

func runDoSomething(cmd *cobra.Command, args []string) error {
    ctx := context.GetGlobal()

    // Get arguments and flags
    input := args[0]
    option, _ := cmd.Flags().GetString("option")
    force, _ := cmd.Flags().GetBool("force")

    // Check dry-run mode
    if ctx.DryRun {
        ctx.Output.DryRun(fmt.Sprintf("Would do something with %s", input))
        return nil
    }

    // Confirm if needed
    if !force && !ctx.Confirm(fmt.Sprintf("Do something with %s?", input), false) {
        ctx.Output.Info("Cancelled")
        return nil
    }

    // Do the work
    logger.L.Info("doing something", "input", input, "option", option)
    result, err := doWork(cmd.Context(), input, option)
    if err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }

    // Output results
    ctx.Output.Success(fmt.Sprintf("Successfully processed: %s", result))
    return nil
}

func doWork(ctx context.Context, input, option string) (string, error) {
    // Your implementation here
    return input + "-processed", nil
}
```

{%- else %}

```go
// internal/cli/myfeature/myfeature.go
package myfeature

import (
    "fmt"

    "github.com/urfave/cli/v2"
    "github.com/fast-ish/${{values.name}}/internal/context"
    "github.com/fast-ish/${{values.name}}/internal/logger"
)

// Cmd is the root command for myfeature
var Cmd = &cli.Command{
    Name:  "myfeature",
    Usage: "My new feature commands",
    Subcommands: []*cli.Command{
        {
            Name:   "do-something",
            Usage:  "Do something useful",
            Action: runDoSomething,
            Flags: []cli.Flag{
                &cli.StringFlag{
                    Name:    "option",
                    Aliases: []string{"o"},
                    Value:   "default",
                    Usage:   "An option",
                },
                &cli.BoolFlag{
                    Name:    "force",
                    Aliases: []string{"f"},
                    Usage:   "Force operation",
                },
            },
        },
    },
}

func runDoSomething(c *cli.Context) error {
    ctx := context.GetGlobal()

    // Get arguments and flags
    if c.NArg() != 1 {
        return fmt.Errorf("requires exactly 1 argument")
    }
    input := c.Args().Get(0)
    option := c.String("option")
    force := c.Bool("force")

    // Check dry-run mode
    if ctx.DryRun {
        ctx.Output.DryRun(fmt.Sprintf("Would do something with %s", input))
        return nil
    }

    // Confirm if needed
    if !force && !ctx.Confirm(fmt.Sprintf("Do something with %s?", input), false) {
        ctx.Output.Info("Cancelled")
        return nil
    }

    // Do the work
    logger.L.Info("doing something", "input", input, "option", option)
    result, err := doWork(c.Context, input, option)
    if err != nil {
        return fmt.Errorf("operation failed: %w", err)
    }

    // Output results
    ctx.Output.Success(fmt.Sprintf("Successfully processed: %s", result))
    return nil
}
```

{%- endif %}

### 3. Register Command

The command auto-registers if you export it:

```go
// internal/cli/root.go

// Import your new command
import (
    "github.com/fast-ish/${{values.name}}/internal/cli/myfeature"
)

func registerCommands() {
    // Your command is automatically available
    rootCmd.AddCommand(myfeature.Cmd)
}
```

### 4. Test Your Command

```bash
# Help
${{values.name}} myfeature --help

# Execute
${{values.name}} myfeature do-something input-value

# With options
${{values.name}} myfeature do-something --option custom input-value

# Dry run
${{values.name}} --dry-run myfeature do-something input-value
```

---

## Adding New Integrations

### 1. Create Client

```go
// internal/client/myservice.go
package client

import (
    "context"
    "fmt"
    "time"

    "github.com/fast-ish/${{values.name}}/internal/config"
)

type MyServiceClient struct {
    base *BaseClient
    cfg  config.MyServiceConfig
}

func NewMyServiceClient(cfg config.MyServiceConfig) *MyServiceClient {
    return &MyServiceClient{
        base: NewBaseClient(
            cfg.BaseURL,
            map[string]string{
                "Authorization": "Bearer " + cfg.Token,
                "Content-Type":  "application/json",
            },
            time.Duration(cfg.Timeout)*time.Second,
        ),
        cfg: cfg,
    }
}

func (c *MyServiceClient) ListResources(ctx context.Context) ([]Resource, error) {
    var result []Resource
    err := c.base.Get(ctx, "/api/v1/resources", nil, &result)
    if err != nil {
        return nil, fmt.Errorf("failed to list resources: %w", err)
    }
    return result, nil
}

type Resource struct {
    ID     string `json:"id"`
    Name   string `json:"name"`
    Status string `json:"status"`
}
```

### 2. Add Configuration

```go
// internal/config/config.go

type Config struct {
    // ... existing fields
    MyService MyServiceConfig `mapstructure:"myservice"`
}

type MyServiceConfig struct {
    Enabled bool   `mapstructure:"enabled"`
    BaseURL string `mapstructure:"base_url"`
    Token   string `mapstructure:"token"`
    Timeout int    `mapstructure:"timeout"`
}
```

### 3. Add to Context

```go
// internal/context/context.go

type Context struct {
    // ... existing fields

    myServiceOnce sync.Once
    myService     *client.MyServiceClient
}

func (c *Context) MyService() *client.MyServiceClient {
    c.myServiceOnce.Do(func() {
        cfg := c.Config.MyService
        c.myService = client.NewMyServiceClient(cfg)
    })
    return c.myService
}
```

### 4. Create Commands

Follow the [Adding New Commands](#adding-new-commands) section to create commands that use your new client.

### 5. Add Configuration

Add to `config.yaml`:

```yaml
myservice:
  enabled: true
  base_url: https://api.myservice.com
  token: ${MYSERVICE_TOKEN}
  timeout: 30
```

Add to `.env.example`:

```bash
MYSERVICE_TOKEN=your-token-here
```

---

{%- if values.aiProvider != "none" %}

## Adding AI Features

### 1. Add New AI Method

```go
// internal/ai/client.go

func (c *Client) Translate(ctx context.Context, text, targetLang string) (string, error) {
    prompt := fmt.Sprintf("Translate the following text to %s:\n\n%s", targetLang, text)
    return c.Chat(ctx, prompt)
}

func (c *Client) CodeReview(ctx context.Context, code string) (string, error) {
    prompt := fmt.Sprintf(`Review this code for:
- Best practices
- Potential bugs
- Performance issues
- Security concerns

Code:
%s

Provide detailed feedback.`, code)

    return c.Chat(ctx, prompt)
}
```

### 2. Add AI Command

```go
// internal/cli/ai/translate.go

var translateCmd = &cobra.Command{
    Use:   "translate [text]",
    Short: "Translate text using AI",
    Args:  cobra.MinimumNArgs(1),
    RunE:  runTranslate,
}

func init() {
    Cmd.AddCommand(translateCmd)
    translateCmd.Flags().StringP("target", "t", "es", "Target language")
}

func runTranslate(cmd *cobra.Command, args []string) error {
    ctx := context.GetGlobal()
    ai := ctx.AI()

    text := strings.Join(args, " ")
    target, _ := cmd.Flags().GetString("target")

    result, err := ai.Translate(cmd.Context(), text, target)
    if err != nil {
        return fmt.Errorf("translation failed: %w", err)
    }

    ctx.Output.Success(result)
    return nil
}
```

{%- endif %}

---

## Custom Output Formats

### 1. Add Custom Formatter

```go
// internal/output/output.go

func (f *Formatter) Tree(data interface{}) {
    // Custom tree visualization
    switch f.format {
    case "tree":
        f.renderTree(data)
    default:
        f.Data(data) // Fall back to standard formats
    }
}

func (f *Formatter) renderTree(data interface{}) {
    // Your tree rendering logic
}
```

### 2. Register Format

{%- if values.cliFramework == "cobra" %}

```go
// internal/cli/root.go

rootCmd.PersistentFlags().StringP("output", "o", "auto",
    "Output format (auto|json|yaml|table|tree)")
```

{%- endif %}

### 3. Use Custom Format

```bash
${{values.name}} --output tree myfeature list
```

---

## Plugin Development

### 1. Define Plugin Interface

```go
// pkg/plugin/plugin.go
package plugin

import (
    "github.com/spf13/cobra"
)

type Plugin interface {
    Name() string
    Commands() []*cobra.Command
    Init() error
    Cleanup() error
}
```

### 2. Implement Plugin

```go
// plugins/sample/sample.go
package sample

import (
    "github.com/spf13/cobra"
    "github.com/fast-ish/${{values.name}}/pkg/plugin"
)

type SamplePlugin struct{}

func New() plugin.Plugin {
    return &SamplePlugin{}
}

func (p *SamplePlugin) Name() string {
    return "sample"
}

func (p *SamplePlugin) Commands() []*cobra.Command {
    return []*cobra.Command{
        {
            Use:   "sample",
            Short: "Sample plugin command",
            RunE: func(cmd *cobra.Command, args []string) error {
                cmd.Println("Hello from sample plugin!")
                return nil
            },
        },
    }
}

func (p *SamplePlugin) Init() error {
    // Initialize plugin
    return nil
}

func (p *SamplePlugin) Cleanup() error {
    // Cleanup resources
    return nil
}
```

### 3. Load Plugins

```go
// internal/cli/root.go

func loadPlugins() {
    plugins := []plugin.Plugin{
        sample.New(),
    }

    for _, p := range plugins {
        if err := p.Init(); err != nil {
            logger.L.Error("failed to init plugin", "name", p.Name(), "error", err)
            continue
        }

        for _, cmd := range p.Commands() {
            rootCmd.AddCommand(cmd)
        }
    }
}
```

---

## Middleware and Hooks

### 1. Pre-Command Middleware

{%- if values.cliFramework == "cobra" %}

```go
// internal/cli/root.go

func preRunHook(cmd *cobra.Command, args []string) error {
    // Runs before every command

    // Example: Check authentication
    if requiresAuth(cmd) && !isAuthenticated() {
        return fmt.Errorf("authentication required")
    }

    // Example: Rate limiting check
    if !checkRateLimit() {
        return fmt.Errorf("rate limit exceeded")
    }

    return nil
}

rootCmd.PersistentPreRunE = preRunHook
```

{%- endif %}

### 2. Post-Command Hook

```go
func postRunHook(cmd *cobra.Command, args []string) error {
    // Runs after every command

    // Example: Log command execution
    logger.L.Info("command completed",
        "command", cmd.Name(),
        "args", args)

    // Example: Send metrics
    recordCommandExecution(cmd.Name())

    return nil
}

rootCmd.PersistentPostRunE = postRunHook
```

### 3. Custom Validation

```go
func validateInput(cmd *cobra.Command, args []string) error {
    // Custom argument validation
    if len(args) == 0 {
        return fmt.Errorf("requires at least 1 argument")
    }

    if !isValidFormat(args[0]) {
        return fmt.Errorf("invalid format: %s", args[0])
    }

    return nil
}

myCmd.Args = validateInput
```

---

## Testing Extensions

### Unit Tests

```go
func TestMyFeature(t *testing.T) {
    // Setup
    ctx := setupTestContext()

    // Execute
    result, err := myFeature(ctx)

    // Assert
    assert.NoError(t, err)
    assert.Equal(t, "expected", result)
}
```

### Integration Tests

{%- if values.e2eTests %}

```go
func TestMyFeatureE2E(t *testing.T) {
    // Skip in short mode
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    // Execute real command
    cmd := exec.Command("${{values.name}}", "myfeature", "do-something", "test-input")
    output, err := cmd.CombinedOutput()

    // Verify
    assert.NoError(t, err)
    assert.Contains(t, string(output), "expected output")
}
```

{%- endif %}

---

## Best Practices

1. **Follow conventions:** Match existing code style and patterns
2. **Add tests:** Write unit tests for new functionality
3. **Document:** Add godoc comments for exported functions
4. **Error handling:** Always wrap errors with context
5. **Logging:** Use structured logging with appropriate levels
6. **Configuration:** Make features configurable
7. **Backward compatibility:** Don't break existing commands
8. **Dry run support:** Implement `--dry-run` for destructive operations

---

## Examples

See the existing integrations for reference:
{%- if values.aiProvider != "none" %}
- `internal/cli/ai/` - AI command implementation
{%- endif %}
{%- for integration in values.integrations %}
- `internal/cli/{{integration}}/` - {{integration|title}} commands
{%- endfor %}

## Need Help?

- Check [Architecture](architecture.md) for system design
- Review [Common Patterns](PATTERNS.md) for code examples
- Open an issue on [GitHub](https://github.com/fast-ish/${{values.name}}/issues)
