# Salesforce OAuth2 CLI Authentication Tool

[![Build Status](https://github.com/mr-menno/sfdc-go-auth-cli/workflows/CI/badge.svg)](https://github.com/mr-menno/sfdc-go-auth-cli/actions)
[![Release](https://github.com/mr-menno/sfdc-go-auth-cli/workflows/Build%20and%20Release/badge.svg)](https://github.com/mr-menno/sfdc-go-auth-cli/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/mr-menno/sfdc-go-auth-cli)](https://goreportcard.com/report/github.com/mr-menno/sfdc-go-auth-cli)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A robust, production-ready Go command-line application that authenticates with Salesforce using OAuth2 and returns access tokens, refresh tokens, and instance URLs in JSON format. Built with modern Go practices and comprehensive CI/CD automation.

## ✨ Features

- ✅ **REQ-1**: Interactive prompts for Salesforce Client ID and Client Secret
- ✅ **REQ-2**: Local OAuth2 callback server with configurable port
- ✅ **REQ-3**: Clean JSON output with access token, refresh token, and instance URL
- ✅ **REQ-4**: Built with Cobra CLI framework for professional command-line experience
- ✅ **REQ-5**: Comprehensive automated test suite with coverage reporting
- ✅ **REQ-6**: All code passes linting validation (golangci-lint) with zero errors
- ✅ **REQ-7**: Custom Salesforce domain support for organizations with custom domains
- ✅ **BUILD-1**: GitHub Actions CI/CD with automated releases and multi-platform builds
- ✅ **DOCS-1**: Complete documentation and usage examples

### 🚀 Additional Features

- **Multi-platform support**: Linux, macOS, Windows (x64 and ARM64)
- **Custom domain support**: Works with custom Salesforce domains (e.g., company.my.salesforce.com)
- **Docker support**: Containerized deployment option
- **Security focused**: Hidden password input, state parameter validation, HTTPS endpoints
- **Professional CLI**: Help system, flag validation, error handling
- **Development tools**: Makefile, linting, formatting, coverage reports
- **Automated releases**: GitHub Actions with downloadable artifacts

## 📋 Prerequisites

1. **Go 1.19+** installed on your system
2. **Salesforce Connected App** configured with:
   - OAuth settings enabled
   - Callback URL set to: `http://localhost:8080/callback` (or your custom port)
   - Required OAuth scopes: `full` and `refresh_token`
   - **Note**: For custom domains, ensure your Connected App is configured in the correct Salesforce org

## 📦 Installation

### Option 1: Download Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/mr-menno/sfdc-go-auth-cli/releases).

```bash
# Linux x64
wget https://github.com/mr-menno/sfdc-go-auth-cli/releases/latest/download/sfdc-auth-linux-amd64
chmod +x sfdc-auth-linux-amd64
sudo mv sfdc-auth-linux-amd64 /usr/local/bin/sfdc-auth

# macOS x64
wget https://github.com/mr-menno/sfdc-go-auth-cli/releases/latest/download/sfdc-auth-darwin-amd64
chmod +x sfdc-auth-darwin-amd64
sudo mv sfdc-auth-darwin-amd64 /usr/local/bin/sfdc-auth

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/mr-menno/sfdc-go-auth-cli/releases/latest/download/sfdc-auth-windows-amd64.exe" -OutFile "sfdc-auth.exe"
```

### Option 2: Build from Source

```bash
git clone https://github.com/mr-menno/sfdc-go-auth-cli.git
cd sfdc-go-auth-cli
make build
# or
go build -o sfdc-auth main.go
```

### Option 3: Docker

```bash
docker pull ghcr.io/mr-menno/sfdc-go-auth-cli:latest
docker run -it --rm -p 8080:8080 ghcr.io/mr-menno/sfdc-go-auth-cli:latest
```

## Usage

### Basic Usage (Interactive Mode)

1. Run the application:

   ```bash
   ./sfdc-auth
   ```

2. When prompted, enter your Salesforce Connected App credentials:

   - **Client ID**: Your Connected App's Consumer Key
   - **Client Secret**: Your Connected App's Consumer Secret (input is hidden)

### Advanced Usage (CLI Flags)

You can also provide credentials and options via command-line flags:

```bash
# Provide credentials via flags
./sfdc-auth --client-id "your_client_id" --client-secret "your_client_secret"

# Use a different port for the callback server
./sfdc-auth --port 9090

# Use a custom Salesforce domain (for organizations with custom domains)
./sfdc-auth --domain "company.my.salesforce.com"

# Combine multiple options
./sfdc-auth --domain "company.my.salesforce.com" --port 9090 --quiet

# Get help
./sfdc-auth --help
```

### Available Flags

- `-c, --client-id`: Salesforce Client ID (Consumer Key)
- `-s, --client-secret`: Salesforce Client Secret (Consumer Secret)
- `-d, --domain`: Salesforce domain (default: login.salesforce.com)
- `-p, --port`: Port for OAuth callback server (default: 8080)
- `-q, --quiet`: Suppress informational output
- `-h, --help`: Show help information

### Custom Domain Support

For organizations using custom Salesforce domains (My Domain), specify your domain using the `--domain` flag:

```bash
# Standard Salesforce login
./sfdc-auth --domain "login.salesforce.com"  # This is the default

# Custom domain examples
./sfdc-auth --domain "company.my.salesforce.com"
./sfdc-auth --domain "acme.my.salesforce.com"
./sfdc-auth --domain "test.sandbox.my.salesforce.com"  # Sandbox with custom domain
```

**Important Notes for Custom Domains:**

- Ensure your Connected App is configured in the correct Salesforce org
- The domain should be the full domain without `https://` prefix
- Custom domains typically follow the pattern: `[company].my.salesforce.com`
- Sandbox domains may include additional identifiers: `[company].[sandbox].my.salesforce.com`

### Authentication Flow

The application will:

1. Start a local server on the specified port (default: `localhost:8080`)
2. Display an authorization URL (using your specified domain)
3. Open your browser to that URL (or copy/paste it manually)
4. Wait for the OAuth callback from Salesforce

After successful authentication, the application outputs JSON with your tokens:

```json
{
  "access_token": "00D...",
  "refresh_token": "5Aep...",
  "instance_url": "https://your-instance.salesforce.com"
}
```

## Setting up a Salesforce Connected App

1. Log in to your Salesforce org
2. Go to **Setup** → **App Manager**
3. Click **New Connected App**
4. Fill in the basic information
5. Enable OAuth Settings:
   - **Callback URL**: `http://localhost:8080/callback`
   - **Selected OAuth Scopes**: Add `Full access (full)` and `Perform requests at any time (refresh_token, offline_access)`
6. Save and note the **Consumer Key** (Client ID) and **Consumer Secret** (Client Secret)

## Security Notes

- The Client Secret input is hidden for security
- A random state parameter is generated for each OAuth flow to prevent CSRF attacks
- The local server only runs during the authentication process
- Tokens are only displayed in the terminal output

## Error Handling

The application handles various error scenarios:

- Invalid client credentials
- OAuth authorization errors
- Network connectivity issues
- Invalid callback responses

## 🛠️ Development

### Development Setup

```bash
# Clone the repository
git clone https://github.com/mr-menno/sfdc-go-auth-cli.git
cd sfdc-go-auth-cli

# Quick setup (recommended)
./scripts/setup-dev.sh

# OR manual setup
make dev-setup
make deps
```

### Available Make Targets

```bash
make help                # Show all available targets
make build              # Build for current platform
make build-all          # Build for all platforms
make test               # Run tests
make test-coverage      # Run tests with coverage report
make lint               # Run linter (MANDATORY - zero errors required)
make fmt                # Format code
make pre-commit         # Run all quality checks (fmt, test, lint)
make run                # Build and run the application
make docker-build       # Build Docker image
make docker-run         # Build and run Docker container
make release            # Create release archives
make clean              # Clean build artifacts
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests for specific Go versions (requires Docker)
docker run --rm -v "$PWD":/usr/src/app -w /usr/src/app golang:1.19 go test -v ./...
docker run --rm -v "$PWD":/usr/src/app -w /usr/src/app golang:1.20 go test -v ./...
```

### Code Quality

The project enforces strict code quality standards:

- **golangci-lint**: Comprehensive linting with **zero tolerance** for errors
- **go fmt**: Automatic code formatting
- **go vet**: Static analysis for common issues
- **go test -race**: Race condition detection
- **Coverage reporting**: Test coverage tracking
- **Mandatory linting**: All code changes must pass linting before completion

#### Linting Requirements

**🚨 CRITICAL**: All code must pass `golangci-lint` validation with **zero errors** before any development work is considered complete.

```bash
# Run linting locally (requires golangci-lint installation)
make lint

# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linting manually
golangci-lint run --timeout=5m

# Optional: Install git pre-commit hook to enforce linting
cp scripts/pre-commit-hook .git/hooks/pre-commit
chmod +x .git/hooks/pre-commit
```

### Project Structure

```
.
├── .github/
│   └── workflows/          # GitHub Actions CI/CD
├── scripts/
│   └── setup-dev.sh       # Development environment setup
├── main.go                 # Main application code
├── main_test.go           # Test suite
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── Dockerfile             # Docker container definition
├── Makefile              # Build automation
├── .gitignore            # Git ignore rules
├── LICENSE               # MIT license
└── README.md             # This file
```

## 🚀 CI/CD Pipeline

The project includes comprehensive GitHub Actions workflows:

### Continuous Integration (`ci.yml`)

- **Linting**: Code quality checks with golangci-lint
- **Multi-version testing**: Tests on Go 1.19, 1.20, 1.21
- **Cross-platform builds**: Linux, macOS, Windows
- **Coverage reporting**: Automated coverage reports

### Release Pipeline (`release.yml`)

- **Automated releases**: Triggered on version tags (`v*`)
- **Multi-platform binaries**: Linux, macOS, Windows (x64 and ARM64)
- **Compressed archives**: `.tar.gz` for Unix, `.zip` for Windows
- **Docker images**: Published to GitHub Container Registry
- **Checksums**: SHA256 checksums for all artifacts

### Creating a Release

1. Create and push a version tag:

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. GitHub Actions will automatically:
   - Run all tests
   - Build binaries for all platforms
   - Create compressed archives
   - Generate checksums
   - Create a GitHub release with all artifacts
   - Build and push Docker images

## 🐳 Docker Usage

### Pull and Run

```bash
# Pull the latest image
docker pull ghcr.io/mr-menno/sfdc-go-auth-cli:latest

# Run interactively
docker run -it --rm -p 8080:8080 ghcr.io/mr-menno/sfdc-go-auth-cli:latest

# Run with flags
docker run -it --rm -p 9090:9090 ghcr.io/mr-menno/sfdc-go-auth-cli:latest --port 9090 --quiet

# Run with custom domain
docker run -it --rm -p 8080:8080 ghcr.io/mr-menno/sfdc-go-auth-cli:latest --domain "company.my.salesforce.com"
```

### Build Locally

```bash
# Build the image
make docker-build

# Run the container
make docker-run
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. **🚨 MANDATORY**: Run tests (`make test`) - all tests must pass
5. **🚨 MANDATORY**: Run linting (`make lint`) - zero errors required
6. Run formatting (`make fmt`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Quality Gates

All contributions must pass these quality gates:

- ✅ **Tests**: All existing and new tests must pass
- ✅ **Linting**: Zero golangci-lint errors (enforced in CI)
- ✅ **Formatting**: Code must be properly formatted (`go fmt`)
- ✅ **Coverage**: Maintain or improve test coverage

**Note**: The CI pipeline will automatically reject any pull requests that fail linting validation.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/mr-menno/sfdc-go-auth-cli/issues)
- **Discussions**: [GitHub Discussions](https://github.com/mr-menno/sfdc-go-auth-cli/discussions)
- **Documentation**: This README and inline code comments

## 📊 Project Status

- ✅ All core requirements implemented
- ✅ Comprehensive test suite
- ✅ CI/CD pipeline with automated releases
- ✅ Multi-platform support
- ✅ Docker containerization
- ✅ Professional documentation
