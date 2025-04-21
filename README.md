# go-url-shortener

A simple, fast URL shortener service built with Go, DynamoDB, and Redis.

## Features

- Shorten long URLs to easy-to-share short URLs
- Custom short IDs
- Redis caching for fast access
- DynamoDB for persistent storage

## Tech Stack

- Go
- DynamoDB
- Redis
- Gin
- Zap
- Viper


## API Endpoints

### Shorten URL
```
POST /api/v1/urls

Request:
{
  "longUrl": "https://example.com/very/long/url/that/needs/shortening",
  "customShortId": "custom-id",    // Optional
}

Response:
{
  "shortUrl": "http://short.url/abc123",
  "longUrl": "https://example.com/very/long/url/that/needs/shortening",
  "shortId": "abc123",
  "createdAt": "2023-05-01T12:34:56Z",
}
```

### Redirect to Long URL
```
GET /:code
```


### Delete URL
```
DELETE /api/v1/urls/:code
```

## API Documentation

The API is documented using Swagger/OpenAPI. You can access the interactive API documentation by running the application and navigating to:

```
http://localhost:8080/swagger/index.html
```

This provides a user-friendly interface to view all available endpoints, understand the request/response formats, and even test the API directly from your browser.

### Generating API Documentation

The Swagger documentation is generated automatically using swaggo:

```bash
swag init -g cmd/api/main.go -o docs
```
