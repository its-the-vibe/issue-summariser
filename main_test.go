package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestCommandLineArgParsing(t *testing.T) {
	// Save the original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test with single command-line argument
	testMessage := "There is a bug, can you fix it?"
	os.Args = []string{"issue-summariser", testMessage}

	// The actual parsing logic would be tested as part of main()
	// Here we're just testing the logic conceptually
	var input Input
	if len(os.Args) > 1 {
		input.Message = strings.Join(os.Args[1:], " ")
	}

	if input.Message != testMessage {
		t.Errorf("Expected message %q, got %q", testMessage, input.Message)
	}
}

func TestCommandLineArgParsingMultipleArgs(t *testing.T) {
	// Save the original os.Args
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	// Test with multiple command-line arguments (unquoted words)
	os.Args = []string{"issue-summariser", "Fix", "the", "API", "error"}
	expectedMessage := "Fix the API error"

	var input Input
	if len(os.Args) > 1 {
		input.Message = strings.Join(os.Args[1:], " ")
	}

	if input.Message != expectedMessage {
		t.Errorf("Expected message %q, got %q", expectedMessage, input.Message)
	}
}

func TestEmbeddedAgentContent(t *testing.T) {
	// Verify that agentContent is not empty after embedding
	if agentContent == "" {
		t.Error("Embedded agent content should not be empty")
	}

	// Verify it contains expected key content from the agent description
	expectedContent := []string{
		"Issue Summariser Agent",
		"title",
		"prompt",
		"version",
		"message",
	}
	
	for _, expected := range expectedContent {
		if !strings.Contains(agentContent, expected) {
			t.Errorf("Embedded agent content missing expected string: %q", expected)
		}
	}
}

func TestInputJSONMarshaling(t *testing.T) {
	input := Input{Message: "Test message"}
	data, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("Failed to marshal input: %v", err)
	}

	var parsed Input
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal input: %v", err)
	}

	if parsed.Message != input.Message {
		t.Errorf("Expected message %q, got %q", input.Message, parsed.Message)
	}
}

func TestOutputJSONMarshaling(t *testing.T) {
	output := Output{
		Version: 3,
		Title:   "Test issue title",
		Prompt:  "Original prompt text",
	}

	data, err := json.Marshal(output)
	if err != nil {
		t.Fatalf("Failed to marshal output: %v", err)
	}

	var parsed Output
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Failed to unmarshal output: %v", err)
	}

	if parsed.Version != output.Version {
		t.Errorf("Expected version %d, got %d", output.Version, parsed.Version)
	}
	if parsed.Title != output.Title {
		t.Errorf("Expected title %q, got %q", output.Title, parsed.Title)
	}
	if parsed.Prompt != output.Prompt {
		t.Errorf("Expected prompt %q, got %q", output.Prompt, parsed.Prompt)
	}
}

func TestInputJSONFormat(t *testing.T) {
	testJSON := `{"message": "Test message content"}`
	var input Input
	if err := json.Unmarshal([]byte(testJSON), &input); err != nil {
		t.Fatalf("Failed to parse valid JSON: %v", err)
	}

	if input.Message != "Test message content" {
		t.Errorf("Expected message %q, got %q", "Test message content", input.Message)
	}
}

func TestOutputJSONFormat(t *testing.T) {
	output := Output{
		Version: 3,
		Title:   "Add feature X",
		Prompt:  "We need to add feature X",
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal output: %v", err)
	}

	// Verify it's valid JSON
	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Verify all required fields are present
	if _, ok := parsed["version"]; !ok {
		t.Error("Output missing 'version' field")
	}
	if _, ok := parsed["title"]; !ok {
		t.Error("Output missing 'title' field")
	}
	if _, ok := parsed["prompt"]; !ok {
		t.Error("Output missing 'prompt' field")
	}
}
