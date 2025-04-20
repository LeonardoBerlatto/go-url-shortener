FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o api ./cmd/api

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api .
COPY --from=builder /app/.env .

EXPOSE 8080
CMD ["./api"]
