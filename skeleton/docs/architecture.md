# Architecture

## Overview

${{values.name}} is a modular CLI tool built using modern Go patterns and best practices. The architecture follows a clean, layered design with clear separation of concerns.

## High-Level Architecture

```
┌─────────────────────────────────────────────────┐
│              CLI Layer (cmd/)                    │
│  Entry point, version info, initialization      │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│         Command Layer (internal/cli/)           │
│  Command routing, plugin registry, flags        │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│       Business Logic Layer (internal/)          │
│  Service clients, AI integration, context       │
└────────────────┬────────────────────────────────┘
                 │
┌────────────────▼────────────────────────────────┐
│      Infrastructure Layer (internal/)           │
│  Config, logging, output, error handling        │
└─────────────────────────────────────────────────┘
```

## Core Components

### 1. CLI Framework

{%- if values.cliFramework == "cobra" %}
**Framework:** Cobra with Viper

- Command routing via `cobra.Command` hierarchy
- Configuration management via Viper (multi-format support)
- Persistent flags for global options
- Pre-run hooks for initialization
{%- else %}
**Framework:** urfave/cli v2

- Command routing via `cli.Command` array
- Simple flag management
- Before hooks for initialization
- Minimal dependencies
{%- endif %}

### 2. Modular Plugin System

The CLI uses a plugin registry pattern for extensibility:

```go
// Commands auto-register at initialization
func init() {
    registry.Register("ai", ai.Cmd)
{%- for integration in values.integrations %}
    registry.Register("{{integration}}", {{integration}}.Cmd)
{%- endfor %}
}
```

**Benefits:**
- New modules don't require root command changes
- Conditional compilation based on features
- Easy to add/remove integrations

### 3. Context Management

Global context provides lazy-loaded clients:

```go
type Context struct {
    Config  *config.Config
    Output  *output.Formatter

    // Lazy-loaded clients (sync.Once pattern)
{%- if values.aiProvider != "none" %}
    aiOnce sync.Once
    ai     *ai.Client
{%- endif %}
{%- for integration in values.integrations %}
    {{integration}}Once sync.Once
    {{integration}}     *client.{{integration|title}}Client
{%- endfor %}
}
```

**Benefits:**
- Clients initialized only when needed
- Thread-safe initialization
- Shared configuration and state

### 4. Service Integration Layer

Each integration follows a consistent pattern:

```
internal/client/
├── base.go          # Shared HTTP client with retry/rate limiting
├── aws.go           # AWS SDK wrapper
├── github.go        # GitHub API client
└── ...
```

All clients:
- Share HTTP transport for connection pooling
- Implement retry with exponential backoff
- Support rate limiting
- Use context for cancellation

{%- if values.aiProvider != "none" %}

### 5. AI Integration

The AI layer abstracts multiple providers:

```go
type Client struct {
    cfg      config.AIConfig
    provider string

    // Provider-specific clients
{%- if values.aiProvider == "bedrock" %}
    bedrock  *bedrockruntime.Client
{%- elif values.aiProvider == "openai" %}
    openai   *openai.Client
{%- elif values.aiProvider == "anthropic" %}
    anthropic *anthropic.Client
{%- elif values.aiProvider == "ollama" %}
    ollama    *api.Client
{%- endif %}
}
```

**Features:**
{%- if "chat" in values.aiFeatures %}
- Chat completions with context
{%- endif %}
{%- if "analyze" in values.aiFeatures %}
- Text analysis
{%- endif %}
{%- if "summarize" in values.aiFeatures %}
- Document summarization
{%- endif %}
{%- if "generate" in values.aiFeatures %}
- Content generation
{%- endif %}

{%- endif %}

### 6. Configuration Management

{%- if values.cliFramework == "cobra" %}
Multi-source configuration via Viper:

1. Environment variables (prefix: `{{values.name | upper}}_`)
2. Config file (`~/.{{values.name}}/config.yaml`)
3. Command-line flags
4. Defaults

Priority: flags > env > config file > defaults
{%- else %}
Simple configuration loading:

1. Config file (`~/.{{values.name}}/config.{{values.configFormat}}`)
2. Environment variables
3. Defaults
{%- endif %}

### 7. Output Formatting

{%- if values.outputFormat == "charm" %}
Rich terminal output using Charm libraries:

- **lipgloss**: Styled output with colors
- **bubbles**: Interactive components (spinners, tables)
- **huh**: Forms and prompts

**Output modes:**
- Auto: Terminal-aware formatting
- JSON: Machine-readable output
- YAML: Human-readable structured output
- Table: Formatted tables with borders
{%- else %}
Simple table output:

- **tablewriter**: ASCII tables
- **JSON/YAML**: Structured output

**Output modes:**
- Table: ASCII tables
- JSON: Machine-readable
- YAML: Human-readable
{%- endif %}

{%- if values.loggingLibrary != "none" %}

### 8. Structured Logging

{%- if values.loggingLibrary == "slog" %}
**Library:** log/slog (Go standard library)

- Structured logging with key-value pairs
- Multiple backends (text, JSON)
- Context-aware logging
{%- elif values.loggingLibrary == "zap" %}
**Library:** Uber's Zap

- High-performance structured logging
- Production/development modes
- Sampling and log levels
{%- elif values.loggingLibrary == "zerolog" %}
**Library:** rs/zerolog

- Zero-allocation logging
- JSON-first design
- Beautiful console output
{%- endif %}

**Log levels:** DEBUG, INFO, WARN, ERROR

{%- endif %}

{%- if values.metrics %}

### 9. Metrics Collection

**Library:** Prometheus client

Metrics exposed on port 9090:
- Request counters
- Duration histograms
- Error rates
- Custom business metrics

{%- endif %}

{%- if values.tracing %}

### 10. Distributed Tracing

**Library:** OpenTelemetry

- Span creation for operations
- Context propagation
- Export to OTLP endpoint

{%- endif %}

## Data Flow

### Command Execution Flow

```
1. User runs command
   ↓
2. CLI framework parses args/flags
   ↓
3. Pre-run hook:
   - Load configuration
   - Initialize context
   - Setup logging
   ↓
4. Command handler executes:
   - Get client from context (lazy load)
   - Call service API
   - Format output
   ↓
5. Post-run hook (cleanup)
```

### Example: AI Chat Flow

```
User: ${{values.name}} ai chat "Hello"
  ↓
CLI parses command and flags
  ↓
ai.chatCmd.RunE() executes
  ↓
ctx.AI() lazy-loads AI client
  ↓
client.Chat(ctx, prompt) calls provider API
  ↓
Response formatted via ctx.Output
  ↓
Result printed to stdout
```

## Error Handling

Errors are wrapped with context:

```go
if err != nil {
    return fmt.Errorf("failed to fetch data: %w", err)
}
```

**Error types:**
- Network errors (retry automatically)
- API errors (user-facing message)
- Configuration errors (exit early)
- Validation errors (show usage)

## Security Considerations

1. **Credentials:** Never log or print API keys
2. **TLS:** Verify certificates by default
3. **Input validation:** Sanitize user input
4. **Secrets:** Use environment variables or secure stores
5. **Dependencies:** Regular security audits

## Performance Optimizations

1. **Connection pooling:** Shared HTTP transport
2. **Lazy loading:** Initialize clients on demand
3. **Context cancellation:** Respect timeouts
4. **Concurrent requests:** Use goroutines where safe
5. **Caching:** Store frequently accessed data

## Extensibility

### Adding New Modules

See [EXTENDING.md](EXTENDING.md) for detailed instructions.

Quick steps:
1. Create command in `internal/cli/<module>/`
2. Create client in `internal/client/<module>.go`
3. Add config section in `internal/config/config.go`
4. Register command in plugin registry

### Adding New AI Providers

1. Implement provider in `internal/ai/client.go`
2. Add configuration in `internal/config/config.go`
3. Update feature flags in template

## Testing Strategy

{%- if values.testFramework == "testify" %}
- **Unit tests:** testify/assert and testify/mock
{%- elif values.testFramework == "ginkgo" %}
- **Unit tests:** Ginkgo BDD framework with Gomega matchers
{%- else %}
- **Unit tests:** Standard Go testing package
{%- endif %}
{%- if values.includeMocks %}
- **Mocking:** mockery-generated mocks
{%- endif %}
{%- if values.e2eTests %}
- **E2E tests:** Integration tests with real services
{%- endif %}

## Deployment

{%- if values.dockerize %}
### Docker

Multi-stage build for minimal image size:
- Builder stage: Compile binary
- Runtime stage: Distroless or Alpine base
{%- endif %}

{%- if values.releaseStrategy == "goreleaser" %}
### GoReleaser

Automated releases with:
- Cross-platform binaries
{%- if values.dockerize %}
- Docker images
{%- endif %}
{%- if values.homebrew %}
- Homebrew tap
{%- endif %}
- GitHub releases
{%- endif %}

## References

- [Getting Started](GETTING_STARTED.md)
- [Common Patterns](PATTERNS.md)
- [Extending the CLI](EXTENDING.md)
- [Troubleshooting](TROUBLESHOOTING.md)
