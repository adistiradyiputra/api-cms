#!/bin/bash

# Test script for Chat API with OPA.co.id Integration
# Make sure to replace YOUR_OPA_TOKEN with a valid token

BASE_URL="http://localhost:8080/api"
TOKEN="YOUR_OPA_TOKEN"

echo "Testing Chat API with OPA.co.id Integration"
echo "=========================================="

# Test 1: Send Message
echo -e "\n1. Testing Send Message endpoint..."
curl -X POST "$BASE_URL/chat/send" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Hello, this is a test message",
    "conversation_id": "",
    "persona": "Normal",
    "response_mode": "short",
    "reranker": "false",
    "model_name": "llama-4"
  }' \
  -w "\nHTTP Status: %{http_code}\n"

# Test 2: Save Chat
echo -e "\n2. Testing Save Chat endpoint..."
curl -X POST "$BASE_URL/chat/save" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "Test user message",
    "stream_message": "Test assistant response",
    "conversation_id": "test_conv_123",
    "api_conversation_id": "test_api_conv_456"
  }' \
  -w "\nHTTP Status: %{http_code}\n"

# Test 3: Delete Conversation
echo -e "\n3. Testing Delete Conversation endpoint..."
curl -X DELETE "$BASE_URL/chat/conversation" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "conversation_id": "test_conv_123"
  }' \
  -w "\nHTTP Status: %{http_code}\n"

# Test 4: Test without token (should fail)
echo -e "\n4. Testing without token (should fail)..."
curl -X POST "$BASE_URL/chat/send" \
  -H "Content-Type: application/json" \
  -d '{
    "content": "This should fail"
  }' \
  -w "\nHTTP Status: %{http_code}\n"

echo -e "\nTest completed!"
echo "Note: For streaming test, use a tool like curl with --no-buffer flag or a web browser"