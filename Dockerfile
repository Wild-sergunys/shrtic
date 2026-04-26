# Dockerfile
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /shrtik ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /shrtic .
COPY --from=builder /app/migrations ./migrations

# Не копируем .env в образ - будет через volumes или env vars

EXPOSE 8080

CMD ["./shrtic"]