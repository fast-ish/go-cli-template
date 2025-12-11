{%- if values.aiProvider != "none" %}
// Package ai provides AI/LLM integration
package ai

import (
	"context"
	"fmt"

	"github.com/fast-ish/${{values.name}}/internal/config"
{%- if values.aiProvider == "bedrock" %}
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
{%- elif values.aiProvider == "openai" %}
	"github.com/sashabaranov/go-openai"
{%- elif values.aiProvider == "anthropic" %}
	"github.com/anthropics/anthropic-sdk-go"
{%- elif values.aiProvider == "ollama" %}
	"github.com/ollama/ollama/api"
{%- endif %}
)

// Client provides AI operations
type Client struct {
	cfg config.AIConfig
{%- if values.aiProvider == "bedrock" %}
	bedrock *bedrockruntime.Client
{%- elif values.aiProvider == "openai" %}
	openai *openai.Client
{%- elif values.aiProvider == "anthropic" %}
	anthropic *anthropic.Client
{%- elif values.aiProvider == "ollama" %}
	ollama *api.Client
{%- endif %}
}

// NewClient creates a new AI client
func NewClient(cfg config.AIConfig) *Client {
	client := &Client{cfg: cfg}
	client.init()
	return client
}

func (c *Client) init() {
{%- if values.aiProvider == "bedrock" %}
	cfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(c.cfg.Region),
	)
	if err != nil {
		panic(fmt.Errorf("failed to load AWS config: %w", err))
	}
	c.bedrock = bedrockruntime.NewFromConfig(cfg)
{%- elif values.aiProvider == "openai" %}
	c.openai = openai.NewClient(c.cfg.APIKey)
{%- elif values.aiProvider == "anthropic" %}
	c.anthropic = anthropic.NewClient(
		anthropic.WithAPIKey(c.cfg.APIKey),
	)
{%- elif values.aiProvider == "ollama" %}
	client, err := api.ClientFromEnvironment()
	if err != nil {
		panic(fmt.Errorf("failed to create Ollama client: %w", err))
	}
	c.ollama = client
{%- endif %}
}

{%- if "chat" in values.aiFeatures %}

// Chat sends a chat message and returns the response
func (c *Client) Chat(ctx context.Context, prompt string) (string, error) {
{%- if values.aiProvider == "bedrock" %}
	input := &bedrockruntime.ConverseInput{
		ModelId: aws.String(c.cfg.Model),
		Messages: []types.Message{
			{
				Role: types.ConversationRoleUser,
				Content: []types.ContentBlock{
					&types.ContentBlockMemberText{
						Value: prompt,
					},
				},
			},
		},
	}

	output, err := c.bedrock.Converse(ctx, input)
	if err != nil {
		return "", fmt.Errorf("bedrock converse failed: %w", err)
	}

	// Extract response text
	if len(output.Output.(*types.ConverseOutputMemberMessage).Value.Content) > 0 {
		if text, ok := output.Output.(*types.ConverseOutputMemberMessage).Value.Content[0].(*types.ContentBlockMemberText); ok {
			return text.Value, nil
		}
	}

	return "", fmt.Errorf("no response from model")
{%- elif values.aiProvider == "openai" %}
	resp, err := c.openai.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model: c.cfg.Model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	})
	if err != nil {
		return "", fmt.Errorf("openai chat failed: %w", err)
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from model")
{%- elif values.aiProvider == "anthropic" %}
	resp, err := c.anthropic.Messages.New(ctx, anthropic.MessageNewParams{
		Model: anthropic.F(c.cfg.Model),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		}),
		MaxTokens: anthropic.Int(1024),
	})
	if err != nil {
		return "", fmt.Errorf("anthropic chat failed: %w", err)
	}

	if len(resp.Content) > 0 {
		if text, ok := resp.Content[0].AsUnion().(anthropic.TextBlock); ok {
			return text.Text, nil
		}
	}

	return "", fmt.Errorf("no response from model")
{%- elif values.aiProvider == "ollama" %}
	req := &api.GenerateRequest{
		Model:  c.cfg.Model,
		Prompt: prompt,
	}

	var response string
	err := c.ollama.Generate(ctx, req, func(resp api.GenerateResponse) error {
		response += resp.Response
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("ollama generate failed: %w", err)
	}

	return response, nil
{%- endif %}
}
{%- endif %}

{%- if "analyze" in values.aiFeatures %}

// Analyze analyzes text and returns insights
func (c *Client) Analyze(ctx context.Context, text string) (string, error) {
	prompt := fmt.Sprintf("Analyze the following text and provide insights:\n\n%s", text)
	return c.Chat(ctx, prompt)
}
{%- endif %}

{%- if "summarize" in values.aiFeatures %}

// Summarize creates a summary of the given text
func (c *Client) Summarize(ctx context.Context, text string) (string, error) {
	prompt := fmt.Sprintf("Summarize the following text concisely:\n\n%s", text)
	return c.Chat(ctx, prompt)
}
{%- endif %}

{%- if "generate" in values.aiFeatures %}

// Generate generates content based on a prompt
func (c *Client) Generate(ctx context.Context, prompt string, options map[string]any) (string, error) {
	// TODO: Support options like temperature, max_tokens, etc.
	return c.Chat(ctx, prompt)
}
{%- endif %}
{%- else %}
// Package ai is a placeholder when AI is not enabled
package ai

import "github.com/fast-ish/${{values.name}}/internal/config"

type Client struct{}

func NewClient(cfg config.AIConfig) *Client {
	return &Client{}
}
{%- endif %}
