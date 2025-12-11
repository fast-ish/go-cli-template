# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of ${{values.name}}
{%- if values.cliFramework == "cobra" %}
- Cobra-based CLI with rich command structure
{%- else %}
- urfave/cli-based CLI framework
{%- endif %}
{%- if values.aiProvider != "none" %}
- AI integration via {{values.aiProvider}}
{%- endif %}
{%- for integration in values.integrations %}
- {{integration|title}} integration
{%- endfor %}
{%- if values.outputFormat == "charm" %}
- Rich terminal output with Charm libraries
{%- endif %}
{%- if values.loggingLibrary != "none" %}
- Structured logging with {{values.loggingLibrary}}
{%- endif %}
{%- if values.metrics %}
- Prometheus metrics support
{%- endif %}
{%- if values.tracing %}
- OpenTelemetry tracing support
{%- endif %}

### Changed

### Deprecated

### Removed

### Fixed

### Security

## [0.1.0] - 2025-12-10

### Added
- Initial scaffolding from go-cli-template

[Unreleased]: https://github.com/fast-ish/${{values.name}}/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/fast-ish/${{values.name}}/releases/tag/v0.1.0
