// LLM client — speaks the OpenAI Chat Completions API which covers
// OpenAI, Anthropic (via proxy), Groq, Mistral, Ollama, LM Studio, etc.
package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/open-gitagent/gitclaw-go/internal/config"
)

// LLMClient communicates with any OpenAI-compatible Chat Completions API.
type LLMClient struct {
	baseURL    string
	apiKey     string
	modelID    string
	provider   string
	httpClient *http.Client
}

// NewLLMClient creates a client for the given model string.
func NewLLMClient(modelStr string) (*LLMClient, error) {
	provider, modelID, baseURL := config.ParseModelString(modelStr)
	if modelID == "" {
		return nil, fmt.Errorf("invalid model string: %q", modelStr)
	}

	if baseURL == "" {
		baseURL = resolveBaseURL(provider)
	}

	apiKey := resolveAPIKey(provider)

	return &LLMClient{
		baseURL:  strings.TrimRight(baseURL, "/"),
		apiKey:   apiKey,
		modelID:  modelID,
		provider: provider,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}, nil
}

// ChatMessage is a message in the conversation.
type ChatMessage struct {
	Role       string     `json:"role"`
	Content    any        `json:"content"` // string or []ContentBlock
	ToolCalls  []LLMToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
}

// LLMToolCall is a tool call from the LLM response.
type LLMToolCall struct {
	ID       string          `json:"id"`
	Type     string          `json:"type"`
	Function LLMToolFunction `json:"function"`
}

// LLMToolFunction is the function name + arguments JSON.
type LLMToolFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// LLMToolDef is a tool definition sent to the API.
type LLMToolDef struct {
	Type     string         `json:"type"`
	Function LLMFunctionDef `json:"function"`
}

// LLMFunctionDef describes a function for the API.
type LLMFunctionDef struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Parameters  map[string]any `json:"parameters"`
}

// ChatRequest is the request body for /chat/completions.
type ChatRequest struct {
	Model       string       `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Tools       []LLMToolDef  `json:"tools,omitempty"`
	Temperature *float64     `json:"temperature,omitempty"`
	MaxTokens   *int         `json:"max_tokens,omitempty"`
	TopP        *float64     `json:"top_p,omitempty"`
}

// ChatResponse is the response from /chat/completions.
type ChatResponse struct {
	ID      string         `json:"id"`
	Choices []ChatChoice   `json:"choices"`
	Usage   *ChatUsage     `json:"usage,omitempty"`
	Model   string         `json:"model"`
}

// ChatChoice is one choice from the response.
type ChatChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// ChatUsage is token usage from the response.
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// Chat sends a non-streaming chat completion request.
func (c *LLMClient) Chat(ctx context.Context, messages []ChatMessage, tools []LLMToolDef, constraints *config.ModelConstraints) (*ChatResponse, error) {
	req := ChatRequest{
		Model:    c.modelID,
		Messages: messages,
	}
	if len(tools) > 0 {
		req.Tools = tools
	}
	if constraints != nil {
		req.Temperature = constraints.Temperature
		req.MaxTokens = constraints.MaxTokens
		req.TopP = constraints.TopP
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling request: %w", err)
	}

	endpoint := c.baseURL + "/chat/completions"
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.apiKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("LLM request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("LLM API error (HTTP %d): %s", resp.StatusCode, truncate(string(respBody), 500))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w (body: %s)", err, truncate(string(respBody), 300))
	}

	return &chatResp, nil
}

// resolveBaseURL returns the default API base URL for a provider.
func resolveBaseURL(provider string) string {
	// Check environment override first
	if url := os.Getenv("GITCLAW_MODEL_BASE_URL"); url != "" {
		return url
	}

	switch provider {
	case "openai":
		return "https://api.openai.com/v1"
	case "anthropic":
		// Anthropic has its own API, but many proxies expose OpenAI-compatible endpoints
		return "https://api.anthropic.com/v1"
	case "groq":
		return "https://api.groq.com/openai/v1"
	case "mistral":
		return "https://api.mistral.ai/v1"
	case "xai":
		return "https://api.x.ai/v1"
	case "deepseek":
		return "https://api.deepseek.com/v1"
	case "openrouter":
		return "https://openrouter.ai/api/v1"
	case "cerebras":
		return "https://api.cerebras.ai/v1"
	case "ollama":
		return "http://localhost:11434/v1"
	case "lmstudio":
		return "http://localhost:1234/v1"
	default:
		return "https://api.openai.com/v1"
	}
}

// resolveAPIKey returns the API key for a provider from environment variables.
func resolveAPIKey(provider string) string {
	envVars := map[string]string{
		"openai":     "OPENAI_API_KEY",
		"anthropic":  "ANTHROPIC_API_KEY",
		"google":     "GEMINI_API_KEY",
		"groq":       "GROQ_API_KEY",
		"xai":        "XAI_API_KEY",
		"mistral":    "MISTRAL_API_KEY",
		"deepseek":   "DEEPSEEK_API_KEY",
		"openrouter": "OPENROUTER_API_KEY",
		"cerebras":   "CEREBRAS_API_KEY",
	}

	if envVar, ok := envVars[provider]; ok {
		return os.Getenv(envVar)
	}
	// Fallback to OPENAI_API_KEY (works for compatible endpoints)
	return os.Getenv("OPENAI_API_KEY")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "…"
}
