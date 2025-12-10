// Package context provides global application context
package context

import (
	"sync"

	"github.com/fast-ish/${{values.name}}/internal/config"
	"github.com/fast-ish/${{values.name}}/internal/output"
{%- if values.aiProvider != "none" %}
	"github.com/fast-ish/${{values.name}}/internal/ai"
{%- endif %}
{%- for integration in values.integrations %}
	"github.com/fast-ish/${{values.name}}/internal/client/{{integration}}"
{%- endfor %}
)

// Context holds the global application state
type Context struct {
	Config  *config.Config
	Output  *output.Formatter
	Verbose int
	DryRun  bool

{%- if values.aiProvider != "none" %}
	// AI client (lazy-loaded)
	aiOnce   sync.Once
	aiClient *ai.Client
{%- endif %}

{%- for integration in values.integrations %}
	// {{integration|title}} client (lazy-loaded)
	{{integration}}Once   sync.Once
	{{integration}}Client *{{integration}}.Client
{%- endfor %}
}

// NewContext creates a new application context
func NewContext(cfg *config.Config) *Context {
	return &Context{
		Config: cfg,
		Output: output.NewFormatter("auto"),
	}
}

{%- if values.aiProvider != "none" %}

// AI returns the AI client (lazy-loaded)
func (c *Context) AI() *ai.Client {
	c.aiOnce.Do(func() {
		c.aiClient = ai.NewClient(c.Config.AI)
	})
	return c.aiClient
}
{%- endif %}

{%- for integration in values.integrations %}

// {{integration|title}} returns the {{integration}} client (lazy-loaded)
func (c *Context) {{integration|title}}() *{{integration}}.Client {
	c.{{integration}}Once.Do(func() {
		c.{{integration}}Client = {{integration}}.NewClient(c.Config.{{integration|title}})
	})
	return c.{{integration}}Client
}
{%- endfor %}

// Confirm asks the user for confirmation
func (c *Context) Confirm(message string, defaultValue bool) bool {
	if c.DryRun {
		return false
	}
	return c.Output.Confirm(message, defaultValue)
}

// Global context singleton
var (
	globalCtx *Context
	globalMu  sync.RWMutex
)

// SetGlobal sets the global context
func SetGlobal(ctx *Context) {
	globalMu.Lock()
	defer globalMu.Unlock()
	globalCtx = ctx
}

// GetGlobal returns the global context
func GetGlobal() *Context {
	globalMu.RLock()
	defer globalMu.RUnlock()
	return globalCtx
}
