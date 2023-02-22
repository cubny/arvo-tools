# Dockerfile.avro-encode

FROM golang:1.16-alpine as builder

# Set the working directory
WORKDIR /app

# Install build dependencies
RUN apk update && \
    apk add --no-cache git

# Copy the source code
COPY encoder .

# Build the binary
RUN go build -o avro_encode

# Use a smaller base image
FROM alpine:3.14

# Copy the binary from the builder image
COPY --from=builder /app/avro_encode /usr/local/bin/

# Set the command to run the binary
CMD ["avro_encode"]
