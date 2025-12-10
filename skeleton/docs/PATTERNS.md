# Common Patterns

This document shows common code patterns used in ${{values.name}} and how to implement them.

## Table of Contents

- [Command Structure](#command-structure)
- [Client Implementation](#client-implementation)
- [Error Handling](#error-handling)
- [Configuration](#configuration)
- [Output Formatting](#output-formatting)
- [Testing](#testing)
- [AI Integration](#ai-integration)

## Command Structure

{%- if values.cliFramework == "cobra" %}

### Basic Cobra Command

```go
package mymodule

import (
    "github.com/spf13/cobra"
    "github.com/fast-ish/${{values.name}}/internal/context"
    "github.com/fast-ish/${{values.name}}/internal/logger"
)

var Cmd = &cobra.Command{
    Use:   "mymodule",
    Short: "My module commands",
    Long:  "Detailed description of what this module does",
}

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List resources",
    RunE:  runList,
}

func init() {
    // Add subcommands
    Cmd.AddCommand(listCmd)

    // Add flags
    listCmd.Flags().StringP("filter", "f", "", "Filter results")
}

func runList(cmd *cobra.Command, args []string) error {
    ctx := context.GetGlobal()

    // Get flag values
    filter, _ := cmd.Flags().GetString("filter")

    // Get client (lazy-loaded)
    client := ctx.MyModule()

    // Call API
    results, err := client.List(cmd.Context(), filter)
    if err != nil {
        return fmt.Errorf("failed to list: %w", err)
    }

    // Format output
    ctx.Output.Data(results)
    return nil
}
```

{%- else %}

### Basic urfave/cli Command

```go
package mymodule

import (
    "github.com/urfave/cli/v2"
    "github.com/fast-ish/${{values.name}}/internal/context"
    "github.com/fast-ish/${{values.name}}/internal/logger"
)

var Cmd = &cli.Command{
    Name:  "mymodule",
    Usage: "My module commands",
    Subcommands: []*cli.Command{
        {
            Name:   "list",
            Usage:  "List resources",
            Action: runList,
            Flags: []cli.Flag{
                &cli.StringFlag{
                    Name:    "filter",
                    Aliases: []string{"f"},
                    Usage:   "Filter results",
                },
            },
        },
    },
}

func runList(c *cli.Context) error {
    ctx := context.GetGlobal()

    // Get flag values
    filter := c.String("filter")

    // Get client (lazy-loaded)
    client := ctx.MyModule()

    // Call API
    results, err := client.List(c.Context, filter)
    if err != nil {
        return fmt.Errorf("failed to list: %w", err)
    }

    // Format output
    ctx.Output.Data(results)
    return nil
}
```

{%- endif %}

## Client Implementation

### Base Client Pattern

```go
package client

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/fast-ish/${{values.name}}/internal/config"
)

type MyModuleClient struct {
    base    *BaseClient
    cfg     config.MyModuleConfig
}

func NewMyModuleClient(cfg config.MyModuleConfig) *MyModuleClient {
    return &MyModuleClient{
        base: NewBaseClient(cfg.BaseURL, map[string]string{
            "Authorization": "Bearer " + cfg.Token,
            "Content-Type":  "application/json",
        }, 30*time.Second),
        cfg: cfg,
    }
}

func (c *MyModuleClient) List(ctx context.Context, filter string) ([]Resource, error) {
    var result []Resource

    query := map[string]string{}
    if filter != "" {
        query["filter"] = filter
    }

    err := c.base.Get(ctx, "/api/resources", query, &result)
    if err != nil {
        return nil, fmt.Errorf("failed to list resources: %w", err)
    }

    return result, nil
}

func (c *MyModuleClient) Get(ctx context.Context, id string) (*Resource, error) {
    var result Resource

    path := fmt.Sprintf("/api/resources/%s", id)
    err := c.base.Get(ctx, path, nil, &result)
    if err != nil {
        return nil, fmt.Errorf("failed to get resource: %w", err)
    }

    return &result, nil
}

func (c *MyModuleClient) Create(ctx context.Context, req CreateRequest) (*Resource, error) {
    var result Resource

    err := c.base.Post(ctx, "/api/resources", req, &result)
    if err != nil {
        return nil, fmt.Errorf("failed to create resource: %w", err)
    }

    return &result, nil
}
```

### HTTP Client with Retry

The `BaseClient` includes retry logic:

```go
// Example usage - retries happen automatically
client := NewMyModuleClient(cfg)
result, err := client.Get(ctx, "resource-id")

// Retries with exponential backoff:
// - Attempt 1: immediate
// - Attempt 2: wait 100ms
// - Attempt 3: wait 200ms
// - Attempt 4: wait 400ms
```

### Rate Limited Client

```go
// Configure rate limiting in BaseClient
client := NewBaseClient(baseURL, headers, timeout)
client.SetRateLimit(10) // 10 requests per second

// Requests automatically throttled
for i := 0; i < 100; i++ {
    client.Get(ctx, "/api/endpoint", nil, &result)
}
```

## Error Handling

### Wrapping Errors

Always wrap errors with context:

```go
result, err := client.GetResource(ctx, id)
if err != nil {
    return fmt.Errorf("failed to get resource %s: %w", id, err)
}
```

### User-Friendly Errors

```go
func (c *MyClient) Get(ctx context.Context, id string) (*Resource, error) {
    var result Resource

    err := c.base.Get(ctx, fmt.Sprintf("/api/%s", id), nil, &result)
    if err != nil {
        // Check for specific HTTP errors
        if httpErr, ok := err.(*HTTPError); ok {
            switch httpErr.StatusCode {
            case 404:
                return nil, fmt.Errorf("resource %s not found", id)
            case 403:
                return nil, fmt.Errorf("access denied to resource %s", id)
            case 429:
                return nil, fmt.Errorf("rate limit exceeded, try again later")
            }
        }
        return nil, fmt.Errorf("failed to get resource: %w", err)
    }

    return &result, nil
}
```

### Logging Errors

```go
{%- if values.loggingLibrary == "slog" %}
logger.L.Error("operation failed",
    "error", err,
    "resource_id", id,
    "operation", "get")
{%- elif values.loggingLibrary == "zap" %}
logger.L.Error("operation failed",
    zap.Error(err),
    zap.String("resource_id", id),
    zap.String("operation", "get"))
{%- elif values.loggingLibrary == "zerolog" %}
logger.L.Error().
    Err(err).
    Str("resource_id", id).
    Str("operation", "get").
    Msg("operation failed")
{%- endif %}
```

## Configuration

### Adding New Config Section

1. **Add to config struct:**

```go
// internal/config/config.go
type Config struct {
    // ... existing fields

    MyModule MyModuleConfig `mapstructure:"mymodule"`
}

type MyModuleConfig struct {
    Enabled bool   `mapstructure:"enabled"`
    BaseURL string `mapstructure:"base_url"`
    Token   string `mapstructure:"token"`
    Timeout int    `mapstructure:"timeout"`
}
```

2. **Set defaults:**

{%- if values.cliFramework == "cobra" %}
```go
func setDefaults() {
    viper.SetDefault("mymodule.enabled", false)
    viper.SetDefault("mymodule.base_url", "https://api.example.com")
    viper.SetDefault("mymodule.timeout", 30)
}
```
{%- else %}
```go
func defaultConfig() *Config {
    return &Config{
        MyModule: MyModuleConfig{
            Enabled: false,
            BaseURL: "https://api.example.com",
            Timeout: 30,
        },
    }
}
```
{%- endif %}

3. **Bind environment variables:**

{%- if values.cliFramework == "cobra" %}
```go
viper.BindEnv("mymodule.token", "{{values.name | upper}}_MYMODULE_TOKEN")
```
{%- else %}
```go
// Environment variables automatically read
// Example: {{values.name | upper}}_MYMODULE_TOKEN
```
{%- endif %}

## Output Formatting

### Structured Data Output

```go
type Resource struct {
    ID        string    `json:"id" yaml:"id"`
    Name      string    `json:"name" yaml:"name"`
    Status    string    `json:"status" yaml:"status"`
    CreatedAt time.Time `json:"created_at" yaml:"created_at"`
}

func runList(cmd *cobra.Command, args []string) error {
    ctx := context.GetGlobal()

    resources, err := ctx.MyModule().List(cmd.Context(), "")
    if err != nil {
        return err
    }

    // Automatically formats based on --output flag
    ctx.Output.Data(resources)
    return nil
}
```

### Custom Table Output

{%- if values.outputFormat == "charm" %}
```go
import "github.com/charmbracelet/bubbles/table"

func displayTable(resources []Resource) {
    columns := []table.Column{
        {Title: "ID", Width: 20},
        {Title: "Name", Width: 30},
        {Title: "Status", Width: 10},
        {Title: "Created", Width: 20},
    }

    rows := make([]table.Row, len(resources))
    for i, r := range resources {
        rows[i] = table.Row{
            r.ID,
            r.Name,
            r.Status,
            r.CreatedAt.Format(time.RFC3339),
        }
    }

    t := table.New(
        table.WithColumns(columns),
        table.WithRows(rows),
        table.WithFocused(true),
    )

    fmt.Println(t.View())
}
```
{%- else %}
```go
import "github.com/olekukonko/tablewriter"

func displayTable(resources []Resource) {
    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"ID", "Name", "Status", "Created"})

    for _, r := range resources {
        table.Append([]string{
            r.ID,
            r.Name,
            r.Status,
            r.CreatedAt.Format(time.RFC3339),
        })
    }

    table.Render()
}
```
{%- endif %}

### Styled Output

{%- if values.outputFormat == "charm" %}
```go
import "github.com/charmbracelet/lipgloss"

var (
    successStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
    errorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
    warningStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
)

func displayStatus(status string) {
    switch status {
    case "success":
        fmt.Println(successStyle.Render("✓ Operation completed"))
    case "error":
        fmt.Println(errorStyle.Render("✗ Operation failed"))
    case "warning":
        fmt.Println(warningStyle.Render("⚠ Operation completed with warnings"))
    }
}
```
{%- endif %}

## Testing

{%- if values.testFramework == "testify" %}

### Unit Test Example

```go
package client

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestMyClient_List(t *testing.T) {
    // Setup
    client := NewMyModuleClient(config.MyModuleConfig{
        BaseURL: "https://api.example.com",
        Token:   "test-token",
    })

    // Execute
    ctx := context.Background()
    results, err := client.List(ctx, "")

    // Assert
    require.NoError(t, err)
    assert.NotEmpty(t, results)
}

func TestMyClient_Get_NotFound(t *testing.T) {
    client := NewMyModuleClient(config.MyModuleConfig{
        BaseURL: "https://api.example.com",
        Token:   "test-token",
    })

    ctx := context.Background()
    result, err := client.Get(ctx, "nonexistent")

    assert.Error(t, err)
    assert.Nil(t, result)
    assert.Contains(t, err.Error(), "not found")
}
```

{%- if values.includeMocks %}

### Using Mocks

```go
package client

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Generate mock: mockery --name=MyModuleClient --output=./mocks
type MyModuleClient interface {
    List(ctx context.Context, filter string) ([]Resource, error)
    Get(ctx context.Context, id string) (*Resource, error)
}

func TestCommandWithMock(t *testing.T) {
    // Create mock
    mockClient := new(mocks.MyModuleClient)

    // Setup expectations
    mockClient.On("List", mock.Anything, "").
        Return([]Resource{{ID: "1", Name: "Test"}}, nil)

    // Use mock
    results, err := mockClient.List(context.Background(), "")

    // Verify
    assert.NoError(t, err)
    assert.Len(t, results, 1)
    mockClient.AssertExpectations(t)
}
```
{%- endif %}

{%- elif values.testFramework == "ginkgo" %}

### BDD Test Example

```go
package client_test

import (
    "context"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    "github.com/fast-ish/${{values.name}}/internal/client"
    "github.com/fast-ish/${{values.name}}/internal/config"
)

var _ = Describe("MyModuleClient", func() {
    var (
        client *client.MyModuleClient
        ctx    context.Context
    )

    BeforeEach(func() {
        ctx = context.Background()
        client = client.NewMyModuleClient(config.MyModuleConfig{
            BaseURL: "https://api.example.com",
            Token:   "test-token",
        })
    })

    Describe("List", func() {
        Context("when filter is empty", func() {
            It("returns all resources", func() {
                results, err := client.List(ctx, "")
                Expect(err).NotTo(HaveOccurred())
                Expect(results).NotTo(BeEmpty())
            })
        })

        Context("when filter is provided", func() {
            It("returns filtered resources", func() {
                results, err := client.List(ctx, "status=active")
                Expect(err).NotTo(HaveOccurred())
                Expect(results).To(HaveLen(1))
            })
        })
    })

    Describe("Get", func() {
        Context("when resource exists", func() {
            It("returns the resource", func() {
                result, err := client.Get(ctx, "test-id")
                Expect(err).NotTo(HaveOccurred())
                Expect(result).NotTo(BeNil())
                Expect(result.ID).To(Equal("test-id"))
            })
        })

        Context("when resource does not exist", func() {
            It("returns an error", func() {
                result, err := client.Get(ctx, "nonexistent")
                Expect(err).To(HaveOccurred())
                Expect(result).To(BeNil())
                Expect(err.Error()).To(ContainSubstring("not found"))
            })
        })
    })
})
```
{%- endif %}

{%- if values.aiProvider != "none" %}

## AI Integration

### Using AI Client

```go
package commands

import (
    "context"
    "fmt"

    "github.com/fast-ish/${{values.name}}/internal/ai"
    "github.com/fast-ish/${{values.name}}/internal/context"
)

func analyzeError(errorLog string) error {
    ctx := context.GetGlobal()
    aiClient := ctx.AI()

    prompt := fmt.Sprintf("Analyze this error log and suggest fixes:\n\n%s", errorLog)

    response, err := aiClient.Chat(context.Background(), prompt)
    if err != nil {
        return fmt.Errorf("AI analysis failed: %w", err)
    }

    ctx.Output.Info(response)
    return nil
}
```

### Structured AI Prompts

```go
func generateDocumentation(code string) (string, error) {
    ctx := context.GetGlobal()
    aiClient := ctx.AI()

    prompt := fmt.Sprintf(`Generate documentation for this Go code.

Code:
%s

Generate:
1. Package description
2. Function descriptions
3. Example usage
4. Edge cases

Format as Markdown.`, code)

    return aiClient.Generate(context.Background(), prompt, ai.GenerateOptions{
        MaxTokens:   2000,
        Temperature: 0.3,
    })
}
```

{%- endif %}

## Context Management

### Adding Lazy-Loaded Client

```go
// internal/context/context.go
type Context struct {
    // ... existing fields

    myModuleOnce sync.Once
    myModule     *client.MyModuleClient
}

func (c *Context) MyModule() *client.MyModuleClient {
    c.myModuleOnce.Do(func() {
        cfg := c.Config.Profile(c.Profile).MyModule
        c.myModule = client.NewMyModuleClient(cfg)
    })
    return c.myModule
}
```

### Sharing Context

```go
// Set global context once
context.SetGlobal(ctx)

// Access from anywhere
func someFunction() {
    ctx := context.GetGlobal()
    client := ctx.MyModule()
}
```

## Best Practices

1. **Always use context:** Pass `context.Context` to all I/O operations
2. **Wrap errors:** Add context to errors with `fmt.Errorf("...: %w", err)`
3. **Log structured data:** Use key-value pairs in logs
4. **Test error paths:** Don't just test happy paths
5. **Use constants:** Define magic strings and numbers as constants
6. **Document exported functions:** Add godoc comments
7. **Keep functions small:** Single responsibility principle
8. **Use interfaces:** Make code testable with interfaces

## References

- [Architecture](architecture.md)
- [Getting Started](GETTING_STARTED.md)
- [Extending the CLI](EXTENDING.md)
