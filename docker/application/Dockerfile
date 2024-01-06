FROM golang:1.21 AS builder
ARG COMPONENT=main

WORKDIR /go/src/app

# Install dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the app
RUN make build-${COMPONENT}

FROM alpine AS runner
ARG COMPONENT=main

WORKDIR /app

# Copy the binary from the build stage
COPY --from=builder /go/src/app/bin/${COMPONENT} ./application

# Run the binary
ENTRYPOINT ["./application"]
