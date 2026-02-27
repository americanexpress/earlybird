# Contributing to EarlyBird

Thank you for your interest in contributing to EarlyBird! This document provides instructions for setting up your development environment, building the project, and running tests.

## Prerequisites

- **Go**: Version 1.23 or higher.
- **Git**: For version control.

## Installation

1.  **Clone the repository**:
    ```bash
    git clone https://github.com/americanexpress/earlybird.git
    cd earlybird
    ```

2.  **Install dependencies**:
    ```bash
    go mod download
    ```

## Building

You can build the project using the provided script or standard Go commands.

### Using the Build Script
The `build.sh` script installs dependencies, runs unit tests, and builds binaries for Linux, Windows, and macOS.
```bash
./build.sh
```
Binaries will be placed in the `binaries/` directory.

### Manual Build
To build a binary for your current local OS:
```bash
go build -o go-earlybird .
```

## Testing

EarlyBird has a comprehensive suite of unit tests.

### Running All Tests
To run all tests in the project:
```bash
go test ./pkg/...
```

### Running Specific Tests
To run tests for a specific package (e.g., the file handling logic):
```bash
go test -v ./pkg/file
```

### Testing GitIgnore Logic
We have extensive tests for `.gitignore` pattern matching. To run these specific tests:
```bash
go test -v ./pkg/file -run Test_ExtendedGitIgnoreSamples
```
This validates against standard GitIgnore templates for Go, Python, Node.js, and Java.

## Project Structure
- `pkg/`: Core library code.
- `pkg/file/`: File handling and ignore logic.
- `pkg/scan/`: Scanning engine and rules.
- `cmd/`: Entry point for the application.
