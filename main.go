package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	copilot "github.com/github/copilot-sdk/go"
)

//go:embed .github/agents/issue-summariser.agent.md
var agentContent string

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
	var input Input

	// Check if a command-line argument is provided
	if len(os.Args) > 1 {
		// Use all command-line arguments joined as the message
		// This handles multi-word messages properly
		input.Message = strings.Join(os.Args[1:], " ")
	} else {
		// Fall back to reading JSON input from stdin
		inputBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			log.Fatalf("Failed to read input: %v", err)
		}

		if err := json.Unmarshal(inputBytes, &input); err != nil {
			log.Fatalf("Failed to parse input JSON: %v", err)
		}
	}

	ctx := context.Background()

	// Create Copilot client
	client := copilot.NewClient(nil)
	if err := client.Start(ctx); err != nil {
		log.Fatalf("Failed to start Copilot client: %v", err)
	}
	defer client.Stop()

	// Create session with the agent description as system message
	session, err := client.CreateSession(ctx, &copilot.SessionConfig{
		Model: "gpt-4.1",
		SystemMessage: &copilot.SystemMessageConfig{
			Mode:    "replace",
			Content: agentContent,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}
	defer session.Destroy()

	// Create the input JSON
	inputJSON, err := json.Marshal(input)
	if err != nil {
		log.Fatalf("Failed to marshal input: %v", err)
	}

	// Send the input message to the agent
	response, err := session.SendAndWait(ctx, copilot.MessageOptions{
		Prompt: string(inputJSON),
	})
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}

	// Parse the response to extract the JSON
	var output Output
	if response.Data.Content != nil {
		content := extractJSON(*response.Data.Content)
		if err := json.Unmarshal([]byte(content), &output); err != nil {
			fmt.Printf("input is: %v\n", string(inputJSON))
			fmt.Printf("response is: %v\n", *response.Data.Content)
			log.Fatalf("Failed to parse response JSON: %v", err)
		}
	} else {
		log.Fatalf("No content in response")
	}

	// Output the result as JSON
	outputBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal output: %v", err)
	}

	fmt.Println(string(outputBytes))
}
