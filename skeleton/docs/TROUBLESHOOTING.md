# Troubleshooting

Common issues and solutions for ${{values.name}}.

## Table of Contents

- [Installation Issues](#installation-issues)
- [Configuration Problems](#configuration-problems)
- [Authentication Errors](#authentication-errors)
- [Network Issues](#network-issues)
- [Performance Problems](#performance-problems)
- [Debug Mode](#debug-mode)

---

## Installation Issues

### Command not found

**Problem:** `${{values.name}}: command not found`

**Solution:**

```bash
# Add Go bin to PATH
export PATH="$PATH:$HOME/go/bin"

# Add to shell profile
echo 'export PATH="$PATH:$HOME/go/bin"' >> ~/.bashrc  # or ~/.zshrc

# Reload shell
source ~/.bashrc  # or source ~/.zshrc
```

### Go version mismatch

**Problem:** `go version go{{values.goVersion}} or higher is required`

**Solution:**

```bash
# Check current version
go version

# Update Go
# Visit https://go.dev/dl/ and download latest version
```

### Build fails

**Problem:** Build errors during `make build`

**Solutions:**

```bash
# Clean and rebuild
make clean
go mod tidy
make build

# Update dependencies
go get -u ./...
go mod tidy

# Clear Go cache
go clean -cache -modcache -i -r
```

---

## Configuration Problems

### Config file not found

**Problem:** `config file not found`

**Solution:**

```bash
# Create config directory
mkdir -p ~/.{{values.name}}

# Copy example config
cp .env.example ~/.{{values.name}}/config.yaml

# Or specify config location
${{values.name}} --config /path/to/config.yaml
```

### Environment variables not loading

**Problem:** Configuration values not being read

**Solution:**

{%- if values.cliFramework == "cobra" %}

```bash
# Check environment variable naming
# Format: {{values.name | upper}}_SECTION_KEY
export {{values.name | upper}}_AI_PROVIDER=bedrock

# Verify config loading
${{values.name}} config show

# Enable debug to see config loading
${{values.name}} --verbose 2 config show
```

{%- else %}

```bash
# Verify environment variables
env | grep {{values.name | upper}}

# Load .env file manually
export $(cat .env | xargs)

# Verify config
${{values.name}} config show
```

{%- endif %}

### Invalid configuration

**Problem:** `invalid configuration: missing required field`

**Solution:**

```bash
# Validate configuration
${{values.name}} config validate

# Show current configuration
${{values.name}} config show

# Check required fields for each integration
${{values.name}} config check
```

---

## Authentication Errors

{%- if values.aiProvider == "bedrock" %}

### AWS Bedrock authentication fails

**Problem:** `failed to load AWS credentials`

**Solutions:**

```bash
# Configure AWS credentials
aws configure

# Or set environment variables
export AWS_REGION=us-west-2
export AWS_PROFILE=default

# Or use AWS SSO
aws sso login --profile your-profile
export AWS_PROFILE=your-profile

# Verify credentials
aws sts get-caller-identity
```

{%- elif values.aiProvider == "openai" %}

### OpenAI API key invalid

**Problem:** `invalid API key`

**Solutions:**

```bash
# Set API key
export OPENAI_API_KEY=sk-...

# Or in config file
echo "ai:
  apiKey: sk-..." >> ~/.{{values.name}}/config.yaml

# Verify API key
curl https://api.openai.com/v1/models \
  -H "Authorization: Bearer $OPENAI_API_KEY"
```

{%- elif values.aiProvider == "anthropic" %}

### Anthropic API key invalid

**Problem:** `authentication failed`

**Solutions:**

```bash
# Set API key
export ANTHROPIC_API_KEY=sk-ant-...

# Or in config file
echo "ai:
  apiKey: sk-ant-..." >> ~/.{{values.name}}/config.yaml

# Verify API key
curl https://api.anthropic.com/v1/messages \
  -H "x-api-key: $ANTHROPIC_API_KEY"
```

{%- elif values.aiProvider == "ollama" %}

### Ollama connection fails

**Problem:** `failed to connect to Ollama`

**Solutions:**

```bash
# Check Ollama is running
ollama list

# Start Ollama service
ollama serve

# Set custom host
export OLLAMA_HOST=http://localhost:11434

# Test connection
curl http://localhost:11434/api/tags
```

{%- endif %}

{%- if "github" in values.integrations %}

### GitHub token invalid

**Problem:** `bad credentials` or `401 Unauthorized`

**Solutions:**

```bash
# Create new token at https://github.com/settings/tokens
# Required scopes: repo, read:org

# Set token
export GITHUB_TOKEN=ghp_...

# Or in config file
echo "github:
  token: ghp_..." >> ~/.{{values.name}}/config.yaml

# Verify token
curl -H "Authorization: token $GITHUB_TOKEN" \
  https://api.github.com/user
```

{%- endif %}

{%- if "slack" in values.integrations %}

### Slack authentication fails

**Problem:** `invalid_auth` or `token_revoked`

**Solutions:**

```bash
# Create new app at https://api.slack.com/apps
# Add scopes: chat:write, channels:read

# Set token
export SLACK_TOKEN=xoxb-...

# Verify token
curl -H "Authorization: Bearer $SLACK_TOKEN" \
  https://slack.com/api/auth.test
```

{%- endif %}

{%- if "kubernetes" in values.integrations %}

### Kubernetes authentication fails

**Problem:** `Unable to connect to cluster`

**Solutions:**

```bash
# Check kubeconfig
kubectl config view

# Set context
kubectl config use-context your-context

# Verify connection
kubectl cluster-info

# Set KUBECONFIG
export KUBECONFIG=~/.kube/config
```

{%- endif %}

---

## Network Issues

### Connection timeout

**Problem:** `context deadline exceeded` or `connection timeout`

**Solutions:**

```bash
# Increase timeout in config
echo "timeout: 60" >> ~/.{{values.name}}/config.yaml

# Check network connectivity
ping api.example.com

# Check proxy settings
echo $HTTP_PROXY
echo $HTTPS_PROXY

# Bypass proxy if needed
unset HTTP_PROXY HTTPS_PROXY
```

### Rate limit exceeded

**Problem:** `rate limit exceeded` or `429 Too Many Requests`

**Solutions:**

```bash
# Wait and retry
sleep 60
${{values.name}} your-command

# Enable rate limiting in config
echo "rateLimitPerSecond: 5" >> ~/.{{values.name}}/config.yaml

# Use exponential backoff
${{values.name}} --retry 3 your-command
```

### SSL/TLS errors

**Problem:** `certificate verify failed` or `TLS handshake timeout`

**Solutions:**

```bash
# Update CA certificates
# macOS
brew update && brew upgrade ca-certificates

# Linux
sudo apt-get update && sudo apt-get install ca-certificates

# Temporary: Skip verification (NOT RECOMMENDED for production)
{%- if "argocd" in values.integrations %}
export ARGOCD_INSECURE=true
{%- endif %}
```

---

## Performance Problems

### Slow command execution

**Problem:** Commands take too long to execute

**Solutions:**

```bash
# Enable connection pooling (already enabled by default)
# Check for network latency
ping api.example.com

# Use parallel execution where possible
${{values.name}} batch-command --parallel 5

# Profile the command
${{values.name}} --verbose 2 your-command 2>&1 | grep "took"
```

### High memory usage

**Problem:** CLI consuming too much memory

**Solutions:**

```bash
# Limit concurrent operations
${{values.name}} bulk-operation --batch-size 10

# Use streaming for large datasets
${{values.name}} export --stream

# Check for memory leaks
go tool pprof http://localhost:6060/debug/pprof/heap
```

---

## Debug Mode

### Enable verbose logging

```bash
# Level 0: Quiet
${{values.name}} --verbose 0 your-command

# Level 1: Info (default)
${{values.name}} --verbose 1 your-command

# Level 2: Debug
${{values.name}} --verbose 2 your-command
```

### JSON logging

```bash
# Set log format to JSON for structured logs
{%- if values.loggingLibrary == "slog" %}
export {{values.name | upper}}_LOG_FORMAT=json
{%- elif values.loggingLibrary == "zap" %}
export {{values.name | upper}}_LOG_FORMAT=json
{%- elif values.loggingLibrary == "zerolog" %}
export {{values.name | upper}}_LOG_FORMAT=json
{%- endif %}

${{values.name}} --verbose 2 your-command
```

### Dry run mode

```bash
# Preview actions without executing
${{values.name}} --dry-run your-command

# Verify what would happen
${{values.name}} --dry-run --verbose 2 your-command
```

### Debug HTTP requests

```bash
# Enable HTTP debug logging
export {{values.name | upper}}_DEBUG_HTTP=true

# See full request/response
${{values.name}} --verbose 2 your-command
```

---

## Common Error Messages

### "no such host"

**Cause:** DNS resolution failed

**Solution:**

```bash
# Check DNS
nslookup api.example.com

# Try different DNS server
export DNS_SERVER=8.8.8.8
```

### "connection refused"

**Cause:** Service not running or firewall blocking

**Solution:**

```bash
# Check service is running
curl http://localhost:port/health

# Check firewall rules
sudo iptables -L

# Check port is not blocked
telnet hostname port
```

### "permission denied"

**Cause:** Insufficient permissions

**Solution:**

```bash
# Check file permissions
ls -la ~/.{{values.name}}/

# Fix permissions
chmod 644 ~/.{{values.name}}/config.yaml
chmod 700 ~/.{{values.name}}/

# Check API token permissions
# Verify token scopes in service dashboard
```

---

## Getting More Help

### Collect diagnostic info

```bash
# Version info
${{values.name}} version

# Configuration
${{values.name}} config show

# Environment
env | grep {{values.name | upper}}

# System info
uname -a
go version
```

### Report an issue

When reporting issues, include:

1. Output of `${{values.name}} version`
2. Full command that failed
3. Complete error message
4. Relevant configuration (redact secrets!)
5. Debug output with `--verbose 2`

### Resources

- [GitHub Issues](https://github.com/fast-ish/${{values.name}}/issues)
- [Documentation](.)
- [Architecture Guide](architecture.md)
- [Common Patterns](PATTERNS.md)

---

## Known Issues

{%- if values.aiProvider == "bedrock" %}
- **Bedrock quota limits:** Some models have low default quotas. Request quota increase in AWS console.
{%- endif %}
{%- if values.aiProvider == "ollama" %}
- **Ollama memory:** Large models require significant RAM. Use smaller models for development.
{%- endif %}
{%- if "kubernetes" in values.integrations %}
- **Kubernetes RBAC:** Ensure service account has required permissions for cluster operations.
{%- endif %}

---

## Still stuck?

If you can't find a solution here:

1. Search [existing issues](https://github.com/fast-ish/${{values.name}}/issues)
2. Check [discussions](https://github.com/fast-ish/${{values.name}}/discussions)
3. Open a [new issue](https://github.com/fast-ish/${{values.name}}/issues/new) with diagnostic info
