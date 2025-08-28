#!/bin/bash

# Test script for LLM plugin - demonstrates all actions

echo "Testing LLM Plugin..."
echo

# Test metadata
echo "=== Metadata ==="
../plugin metadata | jq .
echo

# Test actions
echo "=== Available Actions ==="
../plugin actions | jq .
echo

# Test generate action (requires OPENAI_API_KEY)
echo "=== Generate Test ==="
echo '{"prompt": "Write a haiku about programming", "max_tokens": 100, "temperature": 0.7}' | ../plugin generate | jq .
echo

# Test chat action (requires OPENAI_API_KEY)
echo "=== Chat Test ==="
echo '{"messages": [{"role": "user", "content": "What is Go programming language?"}], "model": "gpt-3.5-turbo"}' | ../plugin chat | jq .
echo

# Test ollama action (requires local Ollama server)
echo "=== Ollama Test ==="
echo '{"prompt": "Explain what a microservice is", "model": "llama2"}' | ../plugin ollama | jq .
echo

echo "Done!"