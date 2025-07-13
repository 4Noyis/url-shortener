# URL Shortener

A high-performance URL shortener service built with Go, featuring REST API endpoints, Bloom filter optimization, and PostgreSQL storage.

## Features

- **URL Shortening**: Convert long URLs into short, shareable links
- **URL Redirection**: Fast 301 redirects with automatic click tracking
- **Bloom Filter Optimization**: Lightning-fast duplicate detection
- **PostgreSQL Storage**: Reliable data persistence with connection pooling
- **REST API**: Clean, RESTful endpoints with proper HTTP status codes
- **Click Analytics**: Track usage statistics for shortened URLs
- **Expiration Support**: Optional URL expiration dates
- **Base62 Encoding**: Generates readable short codes

## Table of Contents

- [Quick Start](#quick-start)
- [API Documentation](#api-documentation)
- [Architecture](#architecture)
- [Installation](#installation)
- [Performance](#performance)

## Quick Start

### Prerequisites

- Go 1.23.2 or higher
- PostgreSQL database
- Environment variables configured

### 1. Clone the repository

```bash
git clone https://github.com/4Noyis/url-shortener.git
cd url-shortener
```

### 2. Install dependencies

```bash
go mod download
```

### 3. Set up environment variables

Create a `.env` file:

```env
DB_USER=your_username
DB_PASSWORD=your_password
DB_HOST=localhost
DB_PORT=5432
DB_NAME=url_shortener
```

### 4. Run database migrations

```bash
psql -d your_database -f migrations/001_create_urls.sql
```

### 5. Start the server

```bash
go run main.go
```

The server will start on `http://localhost:8080`

## API Documentation

### Shorten URL

Create a shortened URL from a long URL.

```http
POST /api/v1/data/shorten
Content-Type: application/json

{
  "long_url": "https://example.com/very/long/url"
}
```

**Response (201 Created):**
```json
{
  "short_url": "n8Z",
  "long_url": "https://example.com/very/long/url",
  "created_at": "2025-07-13T21:15:21.82265Z"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid URL format
- `409 Conflict`: URL already exists
- `500 Internal Server Error`: Server error

### Redirect URL

Redirect to the original URL using the short code.

```http
GET /api/v1/{shortURL}
```

**Response:**
- `301 Moved Permanently`: Successful redirect with `Location` header
- `404 Not Found`: Short URL not found or expired
- `500 Internal Server Error`: Server error

### Example Usage

```bash
# Shorten a URL
curl -X POST http://localhost:8080/api/v1/data/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://github.com"}'

# Response: {"short_url":"voC","long_url":"https://github.com","created_at":"..."}

# Use the shortened URL (browser will redirect)
curl -L http://localhost:8080/api/v1/voC
```

## Architecture

### Clean Architecture Layers

```
┌─────────────────┐
│    Handlers     │  HTTP request/response handling
├─────────────────┤
│    Services     │  Business logic & orchestration
├─────────────────┤
│   Repository    │  Data access layer
├─────────────────┤
│    Database     │  PostgreSQL storage
└─────────────────┘
```

### Key Components

- **Bloom Filter**: Fast duplicate URL detection (O(1) lookup)
- **Base62 Encoding**: Generates short, readable URLs
- **Connection Pooling**: Efficient database connections
- **Click Tracking**: Automatic usage analytics
- **Gin Framework**: High-performance HTTP router

### Project Structure

```
├── main.go                    # Application entry point
├── config/                    # Database configuration
├── internal/
│   ├── dto/                   # Data transfer objects
│   ├── encoding/              # Base62 encoding utilities
│   ├── filter/                # Bloom filter implementation
│   ├── handlers/              # HTTP handlers
│   ├── models/                # Domain models
│   ├── server/                # HTTP server setup
│   ├── service/               # Business logic
│   └── storage/               # Data access layer
└── migrations/                # Database migrations
```

## Installation

### From Source

```bash
# Clone repository
git clone https://github.com/4Noyis/url-shortener.git
cd url-shortener

# Build binary
go build -o url-shortener main.go

# Run
./url-shortener
```


## Performance

### Bloom Filter Optimization

- **Memory Usage**: ~1.2MB for 1M URLs (1% false positive rate)
- **Lookup Speed**: O(1) constant time
- **False Positives**: Only 1% unnecessary database checks
- **False Negatives**: 0% (guaranteed accuracy for "not exists")

### Database Performance

- **Indexed Lookups**: Fast queries on `short_url` and `long_url`
- **Connection Pooling**: Efficient resource usage
- **Atomic Operations**: Thread-safe click tracking

### Scalability

- **Horizontal Scaling**: Stateless application design
- **Database Scaling**: PostgreSQL read replicas support
- **Caching**: Bloom filter reduces database load
- **Load Balancing**: Multiple instance support

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request


## Future Enhancements

- [ ] Custom short URL aliases
- [ ] Bulk URL shortening
- [ ] Analytics dashboard
- [ ] Rate limiting
- [ ] Redis caching layer
- [ ] Docker containerization
- [ ] Comprehensive test suite
- [ ] API documentation with Swagger
- [ ] User authentication
- [ ] URL preview functionality

## Tech Stack

- **Language**: Go 1.23.2
- **Framework**: Gin HTTP framework
- **Database**: PostgreSQL with pgx driver
- **Optimization**: Bloom filters
- **Architecture**: Clean Architecture pattern
- **Encoding**: Base62 for short URLs

