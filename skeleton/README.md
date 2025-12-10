# ${{values.name}}

> ${{values.description}}

## Features

- 🚀 **Modular Architecture** - Extensible plugin-based design
{%- if values.aiProvider != "none" %}
- 🤖 **AI-Powered** - Integrated with {{values.aiProvider}}
{%- endif %}
- 📊 **Rich Output** - Beautiful terminal UI with {{values.outputFormat}}
- 🔧 **Flexible Configuration** - YAML/TOML/JSON support
- 📝 **Structured Logging** - Production-ready logging with {{values.logging}}
{%- if values.metrics %}
- 📈 **Metrics** - Prometheus metrics endpoint
{%- endif %}
{%- if values.tracing %}
- 🔍 **Tracing** - OpenTelemetry distributed tracing
{%- endif %}

## Installation

### From source

```bash
go install github.com/fast-ish/${{values.name}}@latest
```

{%- if values.homebrew %}

### Homebrew

```bash
brew tap fast-ish/tap
brew install ${{values.name}}
```
{%- endif %}

### Docker

```bash
docker pull ghcr.io/fast-ish/${{values.name}}:latest
docker run --rm ghcr.io/fast-ish/${{values.name}} --help
```

## Quick Start

```bash
# Initialize configuration
${{values.name}} config init

{%- if values.aiProvider != "none" %}
# Try AI features
${{values.name}} ai chat "Hello, how can you help me?"
{%- endif %}

# Get help
${{values.name}} --help
```

## Configuration

Configuration file location: `~/.{{values.name}}/config.yaml`

```yaml
{%- if values.aiProvider != "none" %}
ai:
  provider: {{values.aiProvider}}
{%- if values.aiProvider == "bedrock" %}
  region: us-west-2
  model: anthropic.claude-3-sonnet-20240229-v1:0
{%- elif values.aiProvider == "openai" %}
  api_key: your-api-key
  model: gpt-4
{%- endif %}
{%- endif %}

logging:
  level: info
  format: text

{%- if values.metrics %}
metrics:
  enabled: true
  port: 9090
  path: /metrics
{%- endif %}
```

## Development

```bash
# Install dependencies
make install

# Run tests
make test

# Build
make build

# Run locally
make run ARGS="--help"

# Watch mode
make dev
```

## Architecture

This CLI is built with a modular, plugin-based architecture:

```
${{values.name}}/
├── cmd/           # CLI entry point
├── internal/
│   ├── cli/       # Command implementations
│   ├── client/    # Service integrations
│   ├── config/    # Configuration
│   ├── context/   # Global state
│   ├── logger/    # Logging
│   └── output/    # Terminal output
└── pkg/           # Public packages
```

### Extending with Modules

Add new commands by creating a module in `internal/cli/`:

```go
package mymodule

import "github.com/spf13/cobra"

var Cmd = &cobra.Command{
    Use:   "mymodule",
    Short: "My custom module",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}
```

Register in `internal/cli/root.go`:

```go
import "github.com/fast-ish/${{values.name}}/internal/cli/mymodule"

func registerCommands() {
    rootCmd.AddCommand(mymodule.Cmd)
}
```

## License

Proprietary - All Rights Reserved

## Support

For issues and questions, please open an issue on [GitHub](https://github.com/fast-ish/${{values.name}}/issues).
