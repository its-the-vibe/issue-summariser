package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	copilot "github.com/github/copilot-sdk/go"
	"google.golang.org/genai"
)

//go:embed .github/agents/issue-summariser.agent.md
var agentContent string

// AIProvider defines the interface for different AI backends
type AIProvider interface {
	GenerateSummary(ctx context.Context, input Input) (*Output, error)
}

// CopilotProvider implements AIProvider using GitHub Copilot
type CopilotProvider struct {
	Model string
}

// GeminiProvider implements AIProvider using Google Gemini
type GeminiProvider struct {
	APIKey string
	Model  string
}

func (p *GeminiProvider) GenerateSummary(ctx context.Context, input Input) (*Output, error) {
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	generateContentConfig := &genai.GenerateContentConfig{
		SystemInstruction: genai.NewContentFromText(agentContent, genai.RoleUser),
	}

	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	resp, err := client.Models.GenerateContent(ctx, p.Model, genai.Text(string(inputJSON)), generateContentConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	if len(resp.Candidates) == 0 || resp.Candidates[0].Content == nil {
		return nil, errors.New("no content in Gemini response")
	}

	var content strings.Builder
	for _, part := range resp.Candidates[0].Content.Parts {
		content.WriteString(part.Text)
	}

	var output Output
	extractedJSON := extractJSON(content.String())
	if err := json.Unmarshal([]byte(extractedJSON), &output); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response JSON: %w (content: %s)", err, content.String())
	}

	return &output, nil
}

func (p *CopilotProvider) GenerateSummary(ctx context.Context, input Input) (*Output, error) {
	// Create Copilot client
	client := copilot.NewClient(nil)
	if err := client.Start(ctx); err != nil {
		return nil, fmt.Errorf("failed to start Copilot client: %w", err)
	}
	defer client.Stop()

	// Create session with the agent description as system message
	session, err := client.CreateSession(ctx, &copilot.SessionConfig{
		Model: p.Model,
		SystemMessage: &copilot.SystemMessageConfig{
			Mode:    "replace",
			Content: agentContent,
		},
		OnPermissionRequest: copilot.PermissionHandler.ApproveAll,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Destroy()

	// Create the input JSON
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	// Send the input message to the agent
	response, err := session.SendAndWait(ctx, copilot.MessageOptions{
		Prompt: string(inputJSON),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get response: %w", err)
	}

	// Parse the response to extract the JSON
	var output Output
	if d, ok := response.Data.(*copilot.AssistantMessageData); ok {
		content := extractJSON(d.Content)
		if err := json.Unmarshal([]byte(content), &output); err != nil {
			return nil, fmt.Errorf("failed to parse response JSON: %w (content: %s)", err, d.Content)
		}
	} else {
		return nil, errors.New("no content in response")
	}

	return &output, nil
}

// Input represents the JSON input format
type Input struct {
	Message string `json:"message"`
}

// Output represents the JSON output format
type Output struct {
	Version int    `json:"version"`
	Title   string `json:"title"`
	Prompt  string `json:"prompt"`
}

// extractJSON attempts to extract a JSON object from a string that may contain
// surrounding text or markdown code fences (e.g. ```json ... ```).
func extractJSON(s string) string {
	s = strings.TrimSpace(s)
	// Strip markdown code fences if present
	if strings.HasPrefix(s, "```") {
		// Remove the opening fence line and closing fence
		if idx := strings.Index(s, "\n"); idx != -1 {
			s = s[idx+1:]
		}
		if idx := strings.LastIndex(s, "```"); idx != -1 {
			s = s[:idx]
		}
		s = strings.TrimSpace(s)
	}
	// Find the first '{' and match its closing '}' accounting for nesting
	start := strings.Index(s, "{")
	if start == -1 {
		return s
	}
	depth := 0
	inString := false
	escaped := false
	for i := start; i < len(s); i++ {
		ch := s[i]
		if escaped {
			escaped = false
			continue
		}
		if ch == '\\' && inString {
			escaped = true
			continue
		}
		if ch == '"' {
			inString = !inString
			continue
		}
		if inString {
			continue
		}
		if ch == '{' {
			depth++
		} else if ch == '}' {
			depth--
			if depth == 0 {
				return s[start : i+1]
			}
		}
	}
	return s
}

func main() {
	var (
		providerName string
		geminiAPIKey string
		geminiModel  string
		copilotModel string
	)

	flag.StringVar(&providerName, "provider", os.Getenv("ISSUE_SUMMARISER_PROVIDER"), "AI provider to use (copilot or gemini)")
	flag.StringVar(&geminiAPIKey, "gemini-api-key", os.Getenv("GEMINI_API_KEY"), "API key for Gemini (required if provider is gemini)")
	flag.StringVar(&geminiModel, "gemini-model", "gemini-1.5-flash", "Gemini model to use")
	flag.StringVar(&copilotModel, "copilot-model", "gpt-4.1", "Copilot model to use")
	flag.Parse()

	if providerName == "" {
		providerName = "copilot"
	}

	var input Input
	if flag.NArg() > 0 {
		// Use remaining command-line arguments as the message
		input.Message = strings.Join(flag.Args(), " ")
	} else {
		// Fall back to reading JSON input from stdin
		inputBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		// Try parsing as JSON first
		if err := json.Unmarshal(inputBytes, &input); err != nil {
			// If not JSON, treat the whole input as the message
			input.Message = string(inputBytes)
		}
	}

	ctx := context.Background()

	var provider AIProvider
	switch strings.ToLower(providerName) {
	case "gemini":
		if geminiAPIKey == "" {
			log.Fatal("Gemini API key is required. Set GEMINI_API_KEY env var or use --gemini-api-key flag.")
		}
		provider = &GeminiProvider{
			APIKey: geminiAPIKey,
			Model:  geminiModel,
		}
	case "copilot":
		provider = &CopilotProvider{
			Model: copilotModel,
		}
	default:
		log.Fatalf("Unknown provider: %s", providerName)
	}

	output, err := provider.GenerateSummary(ctx, input)
	if err != nil {
		log.Fatalf("Failed to generate summary: %v", err)
	}

	// Output the result as JSON
	outputBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal output: %v", err)
	}

	fmt.Println(string(outputBytes))
}
