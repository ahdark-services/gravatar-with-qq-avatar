FROM golang:1.21 AS builder
ARG COMPONENT=main

WORKDIR /go/src/app

# Install dependencies
RUN apt update && apt install -y \
    make \
    gcc \
    libwebp-dev
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the app
RUN make build-${COMPONENT}

FROM alpine AS runner
ARG COMPONENT=main

WORKDIR /app

# Install dependencies
RUN apk add --no-cache ca-certificates libwebp libwebp-tools

# Copy the binary from the build stage
COPY --from=builder /go/src/app/bin/${COMPONENT} /app/application

# Run the binary
ENTRYPOINT /app/application
