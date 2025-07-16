# Chat API with OPA.co.id Integration

This is a Go Fiber backend application that provides chat functionality with direct OPA.co.id access token validation.

## Features

- **OPA.co.id Authentication**: Direct token validation with OPA.co.id API
- **Chat Functionality**: Send messages, stream responses, and manage conversations
- **File Upload**: Support for image uploads in chat messages
- **Database Storage**: Persistent storage of conversations, messages, and chat history
- **Streaming Support**: Server-Sent Events (SSE) for real-time chat responses

## API Endpoints

### Authentication

All chat endpoints require a valid OPA.co.id access token in the Authorization header:

```
Authorization: Bearer <your_opa_access_token>
```

### Chat Endpoints

#### 1. Send Message
**POST** `/api/chat/send`

Send a message and get a response. Supports both JSON and multipart form data.

**Request Body:**
```json
{
  "content": "Hello, how are you?",
  "conversation_id": "conv_123",
  "persona": "Normal",
  "response_mode": "short",
  "reranker": "false",
  "model_name": "llama-4"
}
```

**File Upload:**
- Use `multipart/form-data` with an `image` field for image uploads
- If only image is sent, content will default to "[Image]"

**Response:**
```json
{
  "status": "success",
  "message": "Pesan berhasil dikirim dan disimpan.",
  "data": {
    "msg": {
      "content": "Hello! I'm doing well, thank you for asking.",
      "role": "assistant"
    },
    "conversation_id": "conv_123",
    "api_conversation_id": "api_conv_456",
    "session_id": "session_789"
  }
}
```

#### 2. Stream Message
**POST** `/api/chat/stream`

Stream chat responses using Server-Sent Events (SSE).

**Request Body:** Same as Send Message

**Response:** SSE stream with events:
```
event: stream
data: {"token": "Hello"}

event: stream
data: {"token": "! "}

event: stream
data: {"token": "I'm"}

event: response
data: {"status": "success", "data": {...}}

event: save_success
data: {"status": "success", "data": {...}}
```

#### 3. Save Chat
**POST** `/api/chat/save`

Save chat data when streaming is handled client-side.

**Request Body:**
```json
{
  "content": "Hello, how are you?",
  "stream_message": "Hello! I'm doing well, thank you for asking.",
  "conversation_id": "conv_123",
  "api_conversation_id": "api_conv_456",
  "image_name": "photo.jpg",
  "image_url": "/uploads/images/uuid.jpg",
  "image_type": "image/jpeg"
}
```

#### 4. Delete Conversation
**DELETE** `/api/chat/conversation`

Delete a conversation (soft delete).

**Request Body:**
```json
{
  "conversation_id": "conv_123"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "Conversation successfully deleted",
  "data": {
    "conversation_id": "conv_123"
  }
}
```

## Database Schema

### Conversations
- `id`: Primary key
- `conversation_id`: Unique conversation identifier
- `api_conversation_id`: External API conversation ID
- `session_id`: User session ID from OPA.co.id
- `user_id`: User ID from OPA.co.id
- `title`: Conversation title (first 100 chars of first message)
- `created_at`, `updated_at`: Timestamps
- `deleted_at`: Soft delete timestamp

### Messages
- `id`: Primary key
- `message_id`: Unique message identifier
- `conversation_id`: Reference to conversation
- `parent_message_id`: Reference to parent message (for threading)
- `role`: "user" or "assistant"
- `content`: Message content
- `message_metadata`: JSON metadata (image info, API response, etc.)
- `created_at`, `updated_at`: Timestamps
- `deleted_at`: Soft delete timestamp

### Chat Histories
- `id`: Primary key
- `chat_id`: Unique chat history identifier
- `conversation_id`: Reference to conversation
- `conversation_session_id`: Session ID for tracking
- `message_user`: User message content
- `message_assistant`: Assistant message content
- `previous_chat_id`: Reference to previous chat
- `file_name`, `file_url`, `file_type`: File attachment info
- `created_at`, `updated_at`: Timestamps
- `deleted_at`: Soft delete timestamp

## Environment Variables

Create a `.env` file with the following variables:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
APP_PORT=8080
```

## Installation and Setup

1. **Install dependencies:**
   ```bash
   go mod tidy
   ```

2. **Set up environment variables:**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Run the application:**
   ```bash
   go run main.go
   ```

4. **Access the API:**
   - Base URL: `http://localhost:8080`
   - API Base: `http://localhost:8080/api`

## OPA.co.id Integration

The application validates access tokens directly with OPA.co.id by making requests to:
- **Endpoint**: `https://api.opa.co.id/validate-token`
- **Method**: POST
- **Headers**: 
  - `Authorization: Bearer <token>`
  - `Content-Type: application/json`

**Expected Response:**
```json
{
  "status": "success",
  "message": "Token valid",
  "data": {
    "user_id": 123,
    "session_id": "session_789",
    "email": "user@example.com",
    "name": "John Doe"
  }
}
```

## File Upload

- **Upload Directory**: `./uploads/images/`
- **URL Path**: `/uploads/images/<filename>`
- **Supported Formats**: All image formats
- **File Size Limit**: 50MB
- **Naming**: UUID-based to prevent conflicts

## Error Handling

All endpoints return consistent error responses:

```json
{
  "error": "Error type",
  "detail": "Detailed error message",
  "message": "User-friendly message"
}
```

Common HTTP status codes:
- `400`: Bad Request (invalid input)
- `401`: Unauthorized (invalid/missing token)
- `404`: Not Found (resource doesn't exist)
- `429`: Too Many Requests (conversation limit reached)
- `500`: Internal Server Error (server error)

## Security Features

- **Token Validation**: All chat endpoints require valid OPA.co.id tokens
- **User Isolation**: Users can only access their own conversations
- **Input Validation**: All inputs are validated and sanitized
- **File Upload Security**: Files are saved with unique names and validated
- **CORS**: Configured for cross-origin requests
- **Rate Limiting**: Conversation limit of 50 messages per conversation