# STAGE 1: Build the binary
FROM golang:1.26.2-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
# We name the binary 'myapp'
RUN CGO_ENABLED=0 GOOS=linux go build -o everlasting .

# STAGE 2: Final minimal image
FROM alpine:latest

WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/everlasting .

# Copy generated docs (if needed)
COPY --from=builder /app/docs ./docs

RUN touch .env

# This sets the executable to run
ENTRYPOINT ["./everlasting"]

# This provides the default argument (dashboard)
# You can override this when running the container
CMD ["dashboard"]