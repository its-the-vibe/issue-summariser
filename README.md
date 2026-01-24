# issue-summariser

A command-line issue summariser agent built with the [GitHub Copilot SDK for Go](https://github.com/github/copilot-sdk/go).

## Description

This tool generates concise, descriptive GitHub issue titles from Slack message content or any other input text. It uses a specialized AI agent that analyzes the message content, extracts the key purpose or problem, and generates a clear issue title following GitHub best practices.

## Prerequisites

Before you begin, make sure you have:

- **Go** 1.25+ installed
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

The issue summariser can be used in two ways:

### Option 1: Command-line Argument (Simplified)

Simply pass your message as a command-line argument:

```bash
./issue-summariser "There is a bug, can you fix it?"
```

This is the recommended and simplest way to use the tool.

### Option 2: JSON Input via stdin (Legacy)

For backward compatibility, the tool still accepts JSON input from stdin:

```bash
echo '{"message": "your message content here"}' | ./issue-summariser
```

### Output Format

Both methods produce the same JSON output:

```json
{
  "version": 4,
  "title": "[generated GitHub issue title]",
  "prompt": "[the exact input message you provided]"
}
```

### Examples

**Example 1: Feature Request**

```bash
./issue-summariser "We need to add support for uploading images to the user profile page. Currently users can only set text-based information but many have requested the ability to upload a profile picture. This should support common formats like PNG, JPG, and GIF."
```

Output:

```json
{
  "version": 4,
  "title": "Add image upload support to user profile page",
  "prompt": "We need to add support for uploading images to the user profile page. Currently users can only set text-based information but many have requested the ability to upload a profile picture. This should support common formats like PNG, JPG, and GIF."
}
```

**Example 2: Bug Report**

```bash
./issue-summariser "The API is returning 500 errors when we try to delete a user that has associated posts. Need to handle this case properly."
```

Expected output:
```json
{
  "version": 4,
  "title": "Fix API error when deleting users with posts",
  "prompt": "The API is returning 500 errors when we try to delete a user that has associated posts. Need to handle this case properly."
}
```

**Example 3: Documentation Update**

```bash
./issue-summariser "Update the documentation to include the new authentication flow we implemented last week"
```

Expected output:
```json
{
  "version": 4,
  "title": "Update documentation for new authentication flow",
  "prompt": "Update the documentation to include the new authentication flow we implemented last week"
}
```

### Legacy JSON Input Examples

For backward compatibility, you can still use JSON input via stdin:

```bash
echo '{"message": "We need to add support for uploading images to the user profile page. Currently users can only set text-based information but many have requested the ability to upload a profile picture. This should support common formats like PNG, JPG, and GIF."}' | ./issue-summariser
```

Output:

```json
{
  "version": 4,
  "title": "Add image upload support to user profile page",
  "prompt": "We need to add support for uploading images to the user profile page. Currently users can only set text-based information but many have requested the ability to upload a profile picture. This should support common formats like PNG, JPG, and GIF."
}
```

### More Legacy Examples

**API Error (JSON Input)**

```bash
echo '{"message": "The API is returning 500 errors when we try to delete a user that has associated posts. Need to handle this case properly."}' | ./issue-summariser
```

Expected output:
```json
{
  "version": 4,
  "title": "Fix API error when deleting users with posts",
  "prompt": "The API is returning 500 errors when we try to delete a user that has associated posts. Need to handle this case properly."
}
```

**Documentation Update (JSON Input)**

```bash
echo '{"message": "Update the documentation to include the new authentication flow we implemented last week"}' | ./issue-summariser
```

Expected output:
```json
{
  "version": 4,
  "title": "Update documentation for new authentication flow",
  "prompt": "Update the documentation to include the new authentication flow we implemented last week"
}
```

## How It Works

The issue summariser uses an embedded AI agent description that contains:

- Instructions for analyzing message content
- Guidelines for generating concise, descriptive issue titles
- Best practices like using imperative mood, being specific, and avoiding vague terms
- Examples to guide the AI agent

The agent uses GitHub Copilot's GPT-4.1 model to analyze the input and generate appropriate issue titles that follow GitHub best practices.

The agent description is embedded directly into the binary using Go's `embed` package, so the tool can run from any location without needing access to the `.github/agents/issue-summariser.agent.md` file.

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
