# ADR 0001: Initial Architecture Decisions

**Status:** Accepted
**Date:** 2025-12-10
**Deciders:** ${{values.owner}}

## Context

We need to build a modular CLI tool that supports AI integration and multiple service integrations. The tool should be extensible, maintainable, and follow Go best practices.

## Decision

### CLI Framework

{%- if values.cliFramework == "cobra" %}
**Selected:** Cobra with Viper

**Rationale:**
- Industry standard for Go CLIs (kubectl, gh, docker)
- Rich command structure with subcommands
- Powerful flag management with persistent flags
- Viper integration for configuration management
- Excellent documentation and community support
{%- else %}
**Selected:** urfave/cli v2

**Rationale:**
- Lightweight and minimal dependencies
- Simple, intuitive API
- Sufficient for our use case
- Faster compilation and smaller binaries
{%- endif %}

### AI Provider

{%- if values.aiProvider == "bedrock" %}
**Selected:** AWS Bedrock

**Rationale:**
- Leverages existing AWS infrastructure
- Multiple model options (Claude, Llama, etc.)
- Enterprise-grade security and compliance
- Pay-per-use pricing
- No separate API key management
{%- elif values.aiProvider == "openai" %}
**Selected:** OpenAI

**Rationale:**
- Industry-leading models (GPT-4, GPT-4-turbo)
- Comprehensive API and SDKs
- Extensive documentation
- Large developer community
{%- elif values.aiProvider == "anthropic" %}
**Selected:** Anthropic Claude

**Rationale:**
- Excellent at code understanding and generation
- Strong safety features
- Large context windows
- Direct API access
{%- elif values.aiProvider == "ollama" %}
**Selected:** Ollama

**Rationale:**
- Local execution (privacy, no API costs)
- Multiple open-source models
- Simple setup and usage
- Good for development and testing
{%- else %}
**Selected:** None

**Rationale:**
- AI features not required for this use case
- Can be added later if needed
{%- endif %}

### Output Formatting

{%- if values.outputFormat == "charm" %}
**Selected:** Charm Libraries (lipgloss, bubbles, huh)

**Rationale:**
- Beautiful terminal UI out of the box
- Interactive components (forms, spinners, tables)
- Consistent styling across commands
- Active development and maintenance
{%- else %}
**Selected:** Simple table output

**Rationale:**
- Minimal dependencies
- Straightforward implementation
- Sufficient for data display
- Easy to pipe to other tools
{%- endif %}

### Logging

{%- if values.loggingLibrary == "slog" %}
**Selected:** log/slog (Go standard library)

**Rationale:**
- Built into Go 1.21+
- Structured logging with zero dependencies
- Multiple backends (text, JSON)
- Future-proof as part of stdlib
{%- elif values.loggingLibrary == "zap" %}
**Selected:** Uber's Zap

**Rationale:**
- High performance
- Rich feature set
- Production-tested at scale
- Excellent structured logging support
{%- elif values.loggingLibrary == "zerolog" %}
**Selected:** rs/zerolog

**Rationale:**
- Zero-allocation logging
- Best performance among logging libraries
- JSON-first design
- Beautiful console output
{%- endif %}

### Configuration Management

{%- if values.cliFramework == "cobra" %}
**Selected:** Viper

**Rationale:**
- Seamless Cobra integration
- Multiple config sources (files, env, flags)
- Config precedence handling
- Hot reload support
- Multiple format support (YAML, JSON, TOML)
{%- else %}
**Selected:** Simple file loading

**Rationale:**
- No additional dependencies
- Direct struct unmarshaling
- Environment variable support
- Sufficient for our needs
{%- endif %}

### Testing Framework

{%- if values.testFramework == "testify" %}
**Selected:** testify

**Rationale:**
- Rich assertion library
- Mock support
- Suite support for complex tests
- Most popular Go testing library
{%- elif values.testFramework == "ginkgo" %}
**Selected:** Ginkgo with Gomega

**Rationale:**
- BDD-style testing
- Excellent test organization
- Rich matchers via Gomega
- Parallel test execution
{%- else %}
**Selected:** Standard Go testing

**Rationale:**
- No dependencies
- Simple and straightforward
- Good enough for basic testing
- Familiar to all Go developers
{%- endif %}

## Architecture Patterns

### Modular Plugin System

Commands are organized as plugins that auto-register at initialization. This allows:
- Easy addition of new commands
- Conditional compilation based on features
- Clear separation of concerns
- Independent testing of modules

### Lazy-Loaded Clients

Service clients are initialized on-demand using sync.Once pattern:
- Reduces startup time
- Saves resources when features aren't used
- Thread-safe initialization
- Shared configuration

### Shared HTTP Transport

All HTTP clients share a connection pool:
- Better resource utilization
- Improved performance
- Consistent retry and timeout behavior
- Easier monitoring and debugging

### Context-Based State

Global context provides shared state:
- Configuration
- Output formatting
- Logging
- Service clients

## Consequences

### Positive

- Clear, maintainable architecture
- Easy to extend with new features
- Good performance characteristics
- Consistent user experience
- Well-tested components

### Negative

- Learning curve for new contributors
- Some abstraction overhead
- Multiple dependencies to manage
- Initial setup complexity

## Alternatives Considered

{%- if values.cliFramework == "cobra" %}
- **urfave/cli**: Simpler but less powerful
- **kingpin**: Good but less active development
- **pflag + custom**: More control but more work
{%- else %}
- **cobra**: More features but heavier
- **kingpin**: Different API style
- **pflag + custom**: Too low-level
{%- endif %}

## References

- [Go CLI Best Practices](https://github.com/golang-standards/project-layout)
{%- if values.cliFramework == "cobra" %}
- [Cobra Documentation](https://github.com/spf13/cobra)
- [Viper Documentation](https://github.com/spf13/viper)
{%- endif %}
{%- if values.outputFormat == "charm" %}
- [Charm Libraries](https://charm.sh/)
{%- endif %}
