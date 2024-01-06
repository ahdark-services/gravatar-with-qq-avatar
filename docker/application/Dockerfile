FROM golang:1.21-bullseye AS builder
ARG COMPONENT=main

WORKDIR /go/src/app

# Install dependencies
RUN apt update && apt install -y libwebp-dev
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the app
RUN go build -o bin/application cmd/${COMPONENT}/main.go

FROM bitnami/minideb:bullseye AS runner

WORKDIR /app

# Install dependencies
RUN apt update && apt install -y libwebp-dev ca-certificates

# Copy the binary from the build stage
COPY --from=builder /go/src/app/bin/application /app/application

# Run the binary
ENTRYPOINT /app/application
