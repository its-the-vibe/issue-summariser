# issue-summariser

A command-line issue summariser agent built with the [GitHub Copilot SDK for Go](https://github.com/github/copilot-sdk/go).

## Description

This tool generates concise, descriptive GitHub issue titles from Slack message content or any other input text. It uses a specialized AI agent that analyzes the message content, extracts the key purpose or problem, and generates a clear issue title following GitHub best practices.

## Prerequisites

Before you begin, make sure you have:

- **Go** 1.21+ installed
- **GitHub Copilot CLI** installed and authenticated ([Installation guide](https://docs.github.com/en/copilot/how-tos/set-up/install-copilot-cli))

Verify the Copilot CLI is working:

```bash
copilot --version
```

## Installation

Clone this repository and build the application:

```bash
git clone https://github.com/its-the-vibe/issue-summariser.git
cd issue-summariser
go build -o issue-summariser main.go
```

## Usage

The issue summariser reads JSON input from stdin and outputs a JSON response with the generated issue title.

### Input Format

```json
{
  "message": "[your message content here]"
}
```

### Output Format

```json
{
  "version": 3,
  "title": "[generated GitHub issue title]",
  "prompt": "[the exact input message you provided]"
}
```

### Example

```bash
echo '{"message": "We need to add support for uploading images to the user profile page. Currently users can only set text-based information but many have requested the ability to upload a profile picture. This should support common formats like PNG, JPG, and GIF."}' | ./issue-summariser
```

Output:

```json
{
  "version": 3,
  "title": "Add image upload support to user profile page",
  "prompt": "We need to add support for uploading images to the user profile page. Currently users can only set text-based information but many have requested the ability to upload a profile picture. This should support common formats like PNG, JPG, and GIF."
}
```

### More Examples

**Example 1: API Error**

```bash
echo '{"message": "The API is returning 500 errors when we try to delete a user that has associated posts. Need to handle this case properly."}' | ./issue-summariser
```

Expected output:
```json
{
  "version": 3,
  "title": "Fix API error when deleting users with posts",
  "prompt": "The API is returning 500 errors when we try to delete a user that has associated posts. Need to handle this case properly."
}
```

**Example 2: Documentation Update**

```bash
echo '{"message": "Update the documentation to include the new authentication flow we implemented last week"}' | ./issue-summariser
```

Expected output:
```json
{
  "version": 3,
  "title": "Update documentation for new authentication flow",
  "prompt": "Update the documentation to include the new authentication flow we implemented last week"
}
```

## How It Works

The issue summariser loads the agent description from `.github/agents/issue-summariser.agent.md`, which contains:

- Instructions for analyzing message content
- Guidelines for generating concise, descriptive issue titles
- Best practices like using imperative mood, being specific, and avoiding vague terms
- Examples to guide the AI agent

The agent uses GitHub Copilot's GPT-4.1 model to analyze the input and generate appropriate issue titles that follow GitHub best practices.

## Title Guidelines

The generated titles follow these best practices:

- **Be concise**: Keep titles between 5-10 words when possible
- **Be specific**: Include the main action or problem
- **Use imperative mood**: Start with a verb (Add, Fix, Update, Create, Remove, etc.)
- **Avoid vague terms**: Don't use "Issue with..." or "Problem about..."
- **Include key context**: Mention the component or area affected
- **No punctuation**: Don't end with a period
- **Capitalize appropriately**: Use title case or sentence case consistently

## License

MIT
