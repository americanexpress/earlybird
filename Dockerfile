## Build EarlyBird binary
FROM golang:1.18-alpine as builder

# CONFIGS
ENV EBVERSION="dev"

# Set the working directory
WORKDIR /app

# Copy the source code
COPY . .

# Download dependencies
RUN apk add gcc
RUN apk add libc-dev
RUN go mod tidy

# Run Tests
RUN go test -p 10 ./pkg/... -covermode=count

# Build the binary
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/americanexpress/earlybird/pkg/buildflags.Version=${EBVERSION}'" -o binaries/go-earlybird-linux 

## Create EarlyBird image
FROM alpine:latest

# Define image tag
LABEL tag=${EBVERSION}

# CONFIGS
ENV ebConfigPath="/root/.go-earlybird"
ENV localBinDir="/usr/local/bin"

# Create required directories
RUN mkdir -p ${ebConfigPath}
RUN mkdir -p ${localBinDir}

# Required files to proper locations
COPY --from=builder /app/config ${ebConfigPath}/
COPY --from=builder /app/.ge_ignore ${ebConfigPath}/.ge_ignore
COPY --from=builder /app/binaries/go-earlybird-linux ${localBinDir}/go-earlybird

# Ensure permissions on file are executable
RUN chmod u+x ${localBinDir}/go-earlybird

# Set the working directory
WORKDIR /app

# Run the binary
CMD ["go-earlybird", "-path=/app/", "-display-severity=high", "-fail-severity=high"]