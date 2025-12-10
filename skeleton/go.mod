module github.com/fast-ish/${{values.name}}

go ${{values.goVersion}}

require (
{%- if values.cliFramework == "cobra" %}
	github.com/spf13/cobra v1.8.1
	github.com/spf13/viper v1.19.0
{%- elif values.cliFramework == "urfave" %}
	github.com/urfave/cli/v2 v2.27.5
{%- endif %}
{%- if values.logging == "slog" %}
	// stdlib slog - no external dependency
{%- elif values.logging == "zap" %}
	go.uber.org/zap v1.27.0
{%- elif values.logging == "zerolog" %}
	github.com/rs/zerolog v1.33.0
{%- endif %}
{%- if values.outputFormat == "charm" %}
	github.com/charmbracelet/lipgloss v1.0.0
	github.com/charmbracelet/bubbles v0.20.0
	github.com/charmbracelet/bubbletea v1.2.4
	github.com/charmbracelet/huh v0.6.0
{%- elif values.outputFormat == "tablewriter" %}
	github.com/olekukonko/tablewriter v0.0.5
{%- endif %}
{%- if values.aiProvider == "bedrock" %}
	github.com/aws/aws-sdk-go-v2 v1.32.6
	github.com/aws/aws-sdk-go-v2/config v1.28.6
	github.com/aws/aws-sdk-go-v2/service/bedrockruntime v1.22.1
{%- elif values.aiProvider == "openai" %}
	github.com/sashabaranov/go-openai v1.35.6
{%- elif values.aiProvider == "anthropic" %}
	github.com/anthropics/anthropic-sdk-go v0.2.0-alpha.5
{%- elif values.aiProvider == "ollama" %}
	github.com/ollama/ollama v0.4.7
{%- endif %}
{%- if "aws" in values.integrations %}
	github.com/aws/aws-sdk-go-v2 v1.32.6
	github.com/aws/aws-sdk-go-v2/config v1.28.6
{%- endif %}
{%- if "github" in values.integrations %}
	github.com/google/go-github/v67 v67.0.0
{%- endif %}
{%- if "kubernetes" in values.integrations %}
	k8s.io/client-go v0.31.3
	k8s.io/api v0.31.3
	k8s.io/apimachinery v0.31.3
{%- endif %}
{%- if values.metrics %}
	github.com/prometheus/client_golang v1.20.5
{%- endif %}
{%- if values.tracing %}
	go.opentelemetry.io/otel v1.32.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.32.0
	go.opentelemetry.io/otel/sdk v1.32.0
{%- endif %}
{%- if values.configFormat == "yaml" or values.configFormat == "all" %}
	gopkg.in/yaml.v3 v3.0.1
{%- endif %}
{%- if values.configFormat == "toml" or values.configFormat == "all" %}
	github.com/pelletier/go-toml/v2 v2.2.3
{%- endif %}
)

require (
{%- if values.testFramework == "testify" %}
	github.com/stretchr/testify v1.10.0
{%- elif values.testFramework == "ginkgo" %}
	github.com/onsi/ginkgo/v2 v2.22.2
	github.com/onsi/gomega v1.36.2
{%- endif %}
{%- if values.includeMocks %}
	github.com/vektra/mockery/v2 v2.50.1
{%- endif %}
)
