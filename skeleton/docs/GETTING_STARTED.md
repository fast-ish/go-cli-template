# Getting Started with ${{values.name}}

This guide will help you get up and running with ${{values.name}}.

## Prerequisites

- Go {{values.goVersion}} or higher
- Git
{%- if "kubernetes" in values.integrations %}
- kubectl (for Kubernetes integration)
{%- endif %}
{%- if "terraform" in values.integrations %}
- Terraform (for Terraform integration)
{%- endif %}
{%- if values.aiProvider == "bedrock" %}
- AWS credentials configured (for Bedrock AI)
{%- elif values.aiProvider == "ollama" %}
- Ollama running locally (for Ollama AI)
{%- endif %}

## Installation

### From Source

```bash
git clone https://github.com/fast-ish/${{values.name}}.git
cd ${{values.name}}
make install
make build
```

This will:
1. Download dependencies
2. Run tests
3. Build the binary to `bin/${{values.name}}`

{%- if values.releaseStrategy == "goreleaser" %}

### Using Go Install

```bash
go install github.com/fast-ish/${{values.name}}/cmd/${{values.name}}@latest
```

{%- if values.homebrew %}

### Using Homebrew

```bash
brew tap fast-ish/tap
brew install ${{values.name}}
```
{%- endif %}

{%- if values.dockerize %}

### Using Docker

```bash
docker pull ghcr.io/fast-ish/${{values.name}}:latest
docker run --rm ghcr.io/fast-ish/${{values.name}}:latest version
```
{%- endif %}
{%- endif %}

## Quick Start

### 1. Verify Installation

```bash
${{values.name}} version
```

Expected output:
```
${{values.name}} version v0.1.0 (built 2025-12-10)
```

### 2. Initialize Configuration

Create a configuration file:

```bash
mkdir -p ~/.{{values.name}}
cp .env.example ~/.{{values.name}}/config.yaml
```

Edit the configuration:

```yaml
{%- if values.aiProvider != "none" %}
# AI Configuration
ai:
{%- if values.aiProvider == "bedrock" %}
  provider: bedrock
  region: us-west-2
  model: anthropic.claude-3-sonnet-20240229-v1:0
{%- elif values.aiProvider == "openai" %}
  provider: openai
  apiKey: ${OPENAI_API_KEY}
  model: gpt-4
{%- elif values.aiProvider == "anthropic" %}
  provider: anthropic
  apiKey: ${ANTHROPIC_API_KEY}
  model: claude-3-sonnet-20240229
{%- elif values.aiProvider == "ollama" %}
  provider: ollama
  host: http://localhost:11434
  model: llama2
{%- endif %}
{%- endif %}

{%- if "github" in values.integrations %}
# GitHub Configuration
github:
  token: ${GITHUB_TOKEN}
  org: fast-ish
{%- endif %}

{%- if "slack" in values.integrations %}
# Slack Configuration
slack:
  token: ${SLACK_TOKEN}
  channel: "#general"
{%- endif %}

# Logging Configuration
logging:
  level: info
  format: text
```

### 3. Set Environment Variables

Create a `.env` file:

```bash
cp .env.example .env
```

Edit and add your credentials:

```bash
{%- if values.aiProvider == "bedrock" %}
AWS_REGION=us-west-2
AWS_PROFILE=default
{%- elif values.aiProvider == "openai" %}
OPENAI_API_KEY=sk-...
{%- elif values.aiProvider == "anthropic" %}
ANTHROPIC_API_KEY=sk-ant-...
{%- endif %}
{%- if "github" in values.integrations %}
GITHUB_TOKEN=ghp_...
{%- endif %}
{%- if "slack" in values.integrations %}
SLACK_TOKEN=xoxb-...
{%- endif %}
```

## Basic Usage

### List Available Commands

```bash
${{values.name}} --help
```

### Common Commands

{%- if values.aiProvider != "none" %}

#### AI Chat

```bash
${{values.name}} ai chat "Explain what this CLI does"
```

{%- if "analyze" in values.aiFeatures %}

#### Analyze Text

```bash
echo "Your text here" | ${{values.name}} ai analyze
```
{%- endif %}

{%- if "summarize" in values.aiFeatures %}

#### Summarize Document

```bash
${{values.name}} ai summarize --file docs/README.md
```
{%- endif %}
{%- endif %}

{%- if "github" in values.integrations %}

#### GitHub Operations

```bash
# List repositories
${{values.name}} github repos list

# Create issue
${{values.name}} github issues create --title "Bug report" --body "Details..."

# List pull requests
${{values.name}} github prs list --state open
```
{%- endif %}

{%- if "kubernetes" in values.integrations %}

#### Kubernetes Operations

```bash
# List pods
${{values.name}} k8s pods list

# Get pod logs
${{values.name}} k8s pods logs my-pod

# Describe deployment
${{values.name}} k8s deployments describe my-deployment
```
{%- endif %}

{%- if "aws" in values.integrations %}

#### AWS Operations

```bash
# List S3 buckets
${{values.name}} aws s3 list

# Get IAM users
${{values.name}} aws iam users list

# Check costs
${{values.name}} aws cost summary --period last-30-days
```
{%- endif %}

### Output Formats

Control output format with `--output` flag:

```bash
# JSON output
${{values.name}} --output json ai chat "Hello"

# YAML output
${{values.name}} --output yaml github repos list

# Table output (default)
${{values.name}} --output table k8s pods list
```

### Verbosity Levels

```bash
# Debug output
${{values.name}} --verbose 2 ai chat "Test"

# Info output (default)
${{values.name}} --verbose 1 ai chat "Test"

# Quiet mode
${{values.name}} --verbose 0 ai chat "Test"
```

### Dry Run Mode

Preview actions without executing:

```bash
${{values.name}} --dry-run github issues create --title "Test"
```

## Configuration Details

### Configuration Precedence

{%- if values.cliFramework == "cobra" %}
Configuration is loaded in this order (highest priority first):

1. Command-line flags
2. Environment variables (prefix: `{{values.name | upper}}_`)
3. Config file (`~/.{{values.name}}/config.yaml`)
4. Defaults
{%- else %}
Configuration is loaded in this order:

1. Config file (`~/.{{values.name}}/config.{{values.configFormat}}`)
2. Environment variables
3. Defaults
{%- endif %}

### Configuration File Locations

The CLI searches for configuration in:

1. `~/.{{values.name}}/config.{{values.configFormat}}`
2. `./config.{{values.configFormat}}`
3. Path specified via `--config` flag

## Development Setup

If you're developing or extending the CLI:

### Install Development Tools

```bash
make deps-dev
```

This installs:
- golangci-lint (linting)
- air (live reload)
{%- if values.includeMocks %}
- mockery (mock generation)
{%- endif %}

### Run in Development Mode

```bash
make dev
```

This starts the CLI with live reload using `air`.

### Run Tests

```bash
# All tests
make test

# With coverage
make test-coverage

# Open coverage report
open coverage.html
```

### Lint Code

```bash
make lint
```

### Format Code

```bash
make format
```

## Common Workflows

### Daily Usage Example

```bash
# Morning standup prep
{%- if "github" in values.integrations %}
${{values.name}} github prs list --author @me --state open
{%- endif %}

# Check system health
{%- if "kubernetes" in values.integrations %}
${{values.name}} k8s pods list --namespace production
{%- endif %}

# Review costs
{%- if "aws" in values.integrations %}
${{values.name}} aws cost summary --period last-7-days
{%- endif %}

# AI-assisted debugging
{%- if values.aiProvider != "none" %}
${{values.name}} ai analyze --file logs/error.log
{%- endif %}
```

## Troubleshooting

### Common Issues

**Command not found:**
```bash
# Add to PATH
export PATH="$PATH:$HOME/go/bin"
```

**Authentication errors:**
```bash
# Verify credentials
${{values.name}} config validate
```

**Configuration not loading:**
```bash
# Check config location
${{values.name}} config show
```

For more troubleshooting tips, see [TROUBLESHOOTING.md](TROUBLESHOOTING.md).

## Next Steps

- Read [Architecture](architecture.md) to understand the system design
- Review [Common Patterns](PATTERNS.md) for code examples
- Learn [How to Extend](EXTENDING.md) the CLI with new features

## Getting Help

- Run `${{values.name}} --help` for command-specific help
- Check the [documentation](.)
- Open an issue on [GitHub](https://github.com/fast-ish/${{values.name}}/issues)
