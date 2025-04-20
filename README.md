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
