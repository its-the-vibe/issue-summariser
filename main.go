package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

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
		// Use command-line argument as the message
		input.Message = os.Args[1]
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

	// Create Copilot client
	client := copilot.NewClient(nil)
	if err := client.Start(); err != nil {
		log.Fatalf("Failed to start Copilot client: %v", err)
	}
	defer client.Stop()

	// Create session with the agent description as system message
	session, err := client.CreateSession(&copilot.SessionConfig{
		Model: "gpt-4.1",
		SystemMessage: &copilot.SystemMessageConfig{
			Content: agentContent,
		},
	})
	if err != nil {
		log.Fatalf("Failed to create session: %v", err)
	}

	// Create the input JSON
	inputJSON, err := json.Marshal(input)
	if err != nil {
		log.Fatalf("Failed to marshal input: %v", err)
	}

	// Send the input message to the agent
	response, err := session.SendAndWait(copilot.MessageOptions{
		Prompt: string(inputJSON),
	}, 0)
	if err != nil {
		log.Fatalf("Failed to get response: %v", err)
	}

	// Parse the response to extract the JSON
	var output Output
	if response.Data.Content != nil {
		if err := json.Unmarshal([]byte(*response.Data.Content), &output); err != nil {
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
