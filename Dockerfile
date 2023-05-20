# Stage 1: Build the Go executable
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o node-connector cmd/main/main.go

# Stage 2: Create a minimal image to run the executable
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/node-connector .

EXPOSE 8080

ENV PORT 8080

CMD ["./node-connector"]
