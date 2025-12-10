# Go CLI Golden Path Template

> The recommended way to build CLI tools at our organization.

[![Backstage](https://img.shields.io/badge/Backstage-Template-blue)](https://backstage.io)
[![Go](https://img.shields.io/badge/Go-1.23-00ADD8)](https://go.dev/)
[![AI](https://img.shields.io/badge/AI-Enabled-purple)](https://aws.amazon.com/bedrock/)
[![License](https://img.shields.io/badge/License-Internal-red)]()

## What's Included

| Category | Features |
|----------|----------|
| **Core** | Cobra/Viper or urfave/cli, modular plugin architecture, lazy-loaded clients |
| **AI Integration** | Bedrock, OpenAI, Anthropic, Ollama - chat, analyze, summarize, generate |
| **Service Integrations** | AWS, GitHub, Kubernetes, Grafana, Slack, Notion, ArgoCD, Terraform |
| **Output** | Rich terminal UI (Charm: lipgloss, bubbles, huh) or simple tables |
| **Observability** | Structured logging (slog, zap, zerolog), Prometheus metrics, OpenTelemetry tracing |
| **Resilience** | HTTP retry with exponential backoff, rate limiting, connection pooling |
| **Security** | Secrets detection, pre-commit hooks, secure defaults, TLS verification |
| **Testing** | testify, Ginkgo/Gomega, mockery, E2E test support |
| **Build & Release** | GoReleaser, cross-platform builds, Docker images, Homebrew tap |
| **DevEx** | Makefile, hot reload (air), VS Code config, golangci-lint |

## Quick Start

### Create a New CLI

1. Go to [Backstage Software Catalog](https://backstage.yourcompany.com/create)
2. Select "Go CLI Tool (Golden Path)"
3. Fill in the form:
   - **Project info**: Name, description, owner
   - **AI provider**: Bedrock (recommended), OpenAI, Anthropic, Ollama, or none
   - **Integrations**: Select services you need to interact with
   - **Features**: Choose output format, logging, testing framework
4. Click "Create"
5. Clone your new repository and start building!

### What You'll Get

```
your-cli/
├── cmd/
│   └── your-cli/           # CLI entry point
├── internal/
│   ├── cli/                # Command definitions
│   │   ├── root.go         # Plugin registry
│   │   ├── ai/             # AI commands (if enabled)
│   │   └── [integration]/  # Integration commands
│   ├── client/             # Service clients
│   ├── config/             # Configuration management
│   ├── context/            # Global state (lazy-loaded)
│   ├── logger/             # Structured logging
│   └── output/             # Output formatting
├── k8s/                    # Kubernetes manifests (if deployed as service)
├── .github/                # CI/CD workflows
│   └── workflows/
│       ├── ci.yaml         # Lint, test, security, build
│       └── release.yaml    # Cross-platform releases
├── docs/                   # Documentation
│   ├── architecture.md     # System design
│   ├── GETTING_STARTED.md  # First steps
│   ├── PATTERNS.md         # Code examples
│   ├── EXTENDING.md        # Customization guide
│   ├── TROUBLESHOOTING.md  # Common issues
│   └── adr/                # Architecture decisions
├── Makefile                # Developer commands
├── Dockerfile              # Multi-stage build
└── README.md               # CLI documentation
```

## Documentation

| Document | Description |
|----------|-------------|
| [Getting Started](./skeleton/docs/GETTING_STARTED.md) | Installation, setup, first commands |
| [Architecture Guide](./skeleton/docs/architecture.md) | System design and components |
| [Patterns Guide](./skeleton/docs/PATTERNS.md) | Code patterns with examples |
| [Extending Guide](./skeleton/docs/EXTENDING.md) | How to add features |
| [Troubleshooting](./skeleton/docs/TROUBLESHOOTING.md) | Common issues and solutions |

## Template Options

### Go Configuration
- **Go Version**: 1.23 (recommended), 1.22, 1.21
- **CLI Framework**:
  - Cobra with Viper (recommended for complex CLIs)
  - urfave/cli v2 (lightweight)

### AI Integration
- **Bedrock**: AWS-hosted models (Claude, Llama, etc.) - recommended for enterprise
- **OpenAI**: GPT-4 for advanced reasoning
- **Anthropic**: Direct Claude API access
- **Ollama**: Local models for privacy/development
- **None**: Skip AI features

**AI Features** (select multiple):
- Chat completions
- Text analysis
- Document summarization
- Content generation

### Service Integrations

Select the services your CLI needs to interact with:

| Integration | Use Cases |
|-------------|-----------|
| **AWS** | S3, IAM, EKS, ECR, Cost Explorer, Bedrock, CloudWatch |
| **GitHub** | Repos, issues, PRs, actions, releases |
| **Kubernetes** | Pods, deployments, nodes, events, logs |
| **Grafana** | Dashboards, alerts, datasources, IRM (incidents, on-call) |
| **Slack** | Messages, channels, notifications |
| **Notion** | Pages, databases, wiki/docs |
| **ArgoCD** | Applications, sync status, health checks |
| **Terraform** | Plan, apply, state management |

### Output & Observability
- **Output Format**: Charm Stack (rich UI) or Simple (tablewriter)
- **Logging**: slog (stdlib), zap (performance), zerolog (zero-allocation)
- **Metrics**: Prometheus (optional)
- **Tracing**: OpenTelemetry (optional)

### Testing & Quality
- **Test Framework**: Standard, testify, or Ginkgo/Gomega
- **Mock Generation**: mockery (optional)
- **E2E Tests**: Integration test setup (optional)

### Build & Release
- **Release Strategy**: GoReleaser (recommended) or Makefile
- **Platforms**: Linux, Darwin (macOS), Windows
- **Docker**: Multi-stage builds with Alpine runtime
- **Homebrew**: Publish to tap (requires GoReleaser)

## Architecture Patterns

Every generated CLI includes these patterns:

| Pattern | Implementation | Benefits |
|---------|----------------|----------|
| **Plugin Registry** | Commands auto-register at init | Easy to add/remove features |
| **Lazy Loading** | Clients initialized on-demand with sync.Once | Fast startup, efficient resources |
| **Shared Transport** | HTTP connection pooling | Better performance, consistent behavior |
| **Context Management** | Global state with lazy clients | Shared config, clean API |
| **Retry & Rate Limiting** | Exponential backoff, token bucket | Resilient against transient failures |
| **Structured Output** | Format-aware (JSON/YAML/Table) | Machine-readable or human-friendly |

## Common Use Cases

### DevOps Automation
```bash
# AI-assisted incident response
my-cli ai analyze --file logs/error.log
my-cli grafana incidents list --status firing

# Kubernetes operations
my-cli k8s pods list --namespace production
my-cli k8s deployments restart api-server

# Cost optimization
my-cli aws cost summary --period last-30-days
my-cli aws cost forecast --days 30
```

### Developer Productivity
```bash
# GitHub workflow automation
my-cli github prs list --author @me --state open
my-cli github issues create --title "Bug" --labels bug,urgent

# AI-powered code review
my-cli ai analyze --file src/main.go
my-cli ai summarize --file docs/design.md

# ArgoCD sync management
my-cli argocd apps sync my-app --prune
my-cli argocd apps status --watch
```

### Platform Operations
```bash
# Multi-service health check
my-cli ops health-check --services api,worker,cache
my-cli grafana dashboards snapshot production-metrics

# Deployment workflows
my-cli deploy rolling --app api --replicas 3
my-cli slack notify "#deploys" "API deployed to production"

# Terraform management
my-cli terraform plan --workspace production
my-cli terraform apply --auto-approve
```

## Architecture Decisions

We've made opinionated choices. Each is documented with rationale:

- [ADR-0001: Initial Architecture](./skeleton/docs/adr/0001-initial-architecture.md)

## Development Workflow

### Local Development

```bash
# Clone your generated CLI
git clone github.com/fast-ish/your-cli
cd your-cli

# Install dependencies
make install

# Run tests
make test

# Build binary
make build

# Run with hot reload
make dev

# Lint code
make lint
```

### Adding a New Command

```bash
# 1. Create command package
mkdir -p internal/cli/myfeature

# 2. Implement command (see docs/EXTENDING.md)
# 3. Test locally
go run ./cmd/your-cli myfeature --help

# 4. Add tests
make test

# 5. Commit and push
git add . && git commit -m "feat: add myfeature command"
git push
```

### Adding a New Integration

See [Extending Guide](./skeleton/docs/EXTENDING.md) for step-by-step instructions on:
- Creating service clients
- Adding configuration
- Implementing commands
- Writing tests

## Best Practices

### Configuration
- Use environment variables for secrets
- Keep sensitive data out of config files
- Use `.env.example` as documentation

### Error Handling
- Always wrap errors with context
- Use structured logging for errors
- Provide user-friendly error messages

### Testing
- Write unit tests for business logic
- Use table-driven tests for multiple cases
- Mock external services

### Security
- Never log or print API keys
- Verify TLS certificates by default
- Run `pre-commit` hooks before committing

## Support

- **Slack**: #platform-help
- **Office Hours**: Thursdays 2-3pm
- **Issues**: Open in this repository
- **Documentation**: See [docs/](./docs)

## Contributing

This template evolves based on team feedback and production learnings.

### Suggesting Changes

1. Open an issue describing the problem or enhancement
2. Discuss with platform team
3. Submit PR with changes
4. Include ADR for significant architectural changes

### What Makes a Good Addition?

- Solves a common problem across teams
- Has been proven in production
- Doesn't add unnecessary complexity
- Has clear documentation and examples
- Maintains backward compatibility
- Follows Go idioms and best practices

### Contributing Code

```bash
# Fork and clone
git clone github.com/fast-ish/go-cli-template
cd go-cli-template

# Create branch
git checkout -b feature/my-enhancement

# Make changes to skeleton/
# Update documentation if needed

# Test the template
# (create a test CLI from Backstage or manually)

# Commit with conventional commits
git commit -m "feat: add support for XYZ integration"

# Push and create PR
git push origin feature/my-enhancement
```

## Examples

### Real CLIs Built with This Template

- **devctl**: Unified DevOps CLI (AWS, GitHub, K8s, Grafana, ArgoCD)
- **ai-assist**: AI-powered developer assistant (code review, analysis, generation)
- **cost-optimizer**: AWS cost analysis and optimization recommendations

## Version History

| Version | Date | Changes |
|---------|------|---------|
| 1.0.0 | 2025-12 | Initial release with AI integration |

## FAQ

**Q: Can I use multiple AI providers?**
A: Each CLI supports one provider at generation time, but you can manually add multiple providers after generation.

**Q: How do I add a new service integration after generation?**
A: See [Extending Guide](./skeleton/docs/EXTENDING.md#adding-new-integrations) for step-by-step instructions.

**Q: Can I deploy a CLI as a Kubernetes service?**
A: Yes! The template includes Kubernetes manifests. This is useful for cron jobs or long-running processes.

**Q: How do I handle secrets?**
A: Use environment variables or AWS Secrets Manager. Never commit secrets to git. The template includes pre-commit hooks with secrets detection.

**Q: What if I need a different CLI framework?**
A: The template supports Cobra and urfave/cli. For other frameworks, you can manually adapt the generated code.

---

🤘 Platform Team
