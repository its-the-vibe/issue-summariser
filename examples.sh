#!/bin/bash

# Example usage script for the issue-summariser CLI
# This script demonstrates how to use the tool with various test inputs

echo "Issue Summariser - Example Usage"
echo "================================="
echo ""

# Example 1: Image upload feature (new command-line argument format)
echo "Example 1: Feature request for image uploads (new format)"
./issue-summariser "We need to add support for uploading images to the user profile page. Currently users can only set text-based information but many have requested the ability to upload a profile picture. This should support common formats like PNG, JPG, and GIF."
echo ""

# Example 2: API error (new command-line argument format)
echo "Example 2: API error when deleting users (new format)"
./issue-summariser "The API is returning 500 errors when we try to delete a user that has associated posts. Need to handle this case properly."
echo ""

# Example 3: Documentation update (new command-line argument format)
echo "Example 3: Documentation update request (new format)"
./issue-summariser "Update the documentation to include the new authentication flow we implemented last week"
echo ""

echo "Legacy JSON format examples:"
echo "----------------------------"
echo ""

# Example 4: Image upload feature (legacy JSON format)
echo "Example 4: Feature request for image uploads (legacy JSON format)"
echo '{"message": "We need to add support for uploading images to the user profile page. Currently users can only set text-based information but many have requested the ability to upload a profile picture. This should support common formats like PNG, JPG, and GIF."}' | ./issue-summariser
echo ""

# Example 5: API error (legacy JSON format)
echo "Example 5: API error when deleting users (legacy JSON format)"
echo '{"message": "The API is returning 500 errors when we try to delete a user that has associated posts. Need to handle this case properly."}' | ./issue-summariser
echo ""

# Example 6: Documentation update (legacy JSON format)
echo "Example 6: Documentation update request (legacy JSON format)"
echo '{"message": "Update the documentation to include the new authentication flow we implemented last week"}' | ./issue-summariser
echo ""

echo "Note: This script requires the GitHub Copilot CLI to be installed and authenticated."
echo "Install from: https://docs.github.com/en/copilot/how-tos/set-up/install-copilot-cli"
