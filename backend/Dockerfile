# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Install git for dependencies if needed
RUN apk add --no-cache git

COPY go.mod ./
COPY . .

RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /bin/server ./cmd/server

# Stage 2: Minimal runtime image
FROM alpine:3.20

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /bin/server /app/server

# Create uploads storage mount point
RUN mkdir -p /data/uploads && chmod 777 /data/uploads

EXPOSE 8080

ENTRYPOINT ["/app/server"]
