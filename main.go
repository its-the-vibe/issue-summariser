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
		if err := json.Unmarshal([]byte(*response.Data.Content), &output); err != nil {
			fmt.Println("input is: %v", string(inputJSON))
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
