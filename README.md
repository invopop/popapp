# Invopop App Template

This is a template repository for building Invopop applications using Go. It provides a clean, modular architecture with web and gateway interfaces, making it easy to create new services that integrate with the Invopop ecosystem.

## 🏗️ Architecture Overview

The template follows a clean architecture pattern with clear separation of concerns:

```
app/
├── main.go                 # Entry point: loads config and runs the CLI
├── app.go                  # Wires the domain and interfaces together
├── config/                 # Configuration files
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── domain/            # Business logic layer
│   │   └── models/        # Domain models
│   └── interfaces/        # External interfaces
│       ├── gateway/       # NATS gateway for async tasks
│       └── web/           # HTTP web interface
│           ├── assets/    # Static assets (embedded)
│           └── components/ # Templ components
└── pkg/                   # Public packages and utilities
```

### Key Components

- **Main Application** (`main.go`): Entry point with CLI commands using Cobra
- **Web Interface**: HTTP server built with Echo framework and Templ templates
- **Gateway Interface**: NATS-based async task processing
- **Domain Layer**: Business logic and models
- **Configuration**: YAML-based config with environment variable support

## 📋 Prerequisites

Before using this template, ensure you have the following dependencies installed:

### Required Dependencies

1. **Go 1.26.0+**
   ```bash
   # Check your Go version
   go version
   ```

2. **Mage** - Build automation tool
   ```bash
   go install github.com/magefile/mage@latest
   ```

3. **Templ** - Type-safe Go templating
   ```bash
   go install github.com/a-h/templ/cmd/templ@latest
   ```

4. **Air** - Live reload for Go apps (development)
   ```bash
   go install github.com/air-verse/air@latest
   ```

### Code Dependencies

- **Echo v4**: HTTP web framework
- **Cobra**: CLI commands and flags
- **Templ**: Type-safe Go templating
- **Zerolog**: Structured logging
- **Invopop Client**: Integration with Invopop services

### Optional Dependencies

5. **Docker** - For containerized development and deployment
6. **NATS Server** - For gateway functionality (can run via Docker)

## 🚀 Quick Start

### 1. Clone and Setup

```bash
# Clone the template
git clone https://github.com/invopop/popapp.git my-new-app
cd my-new-app

# Update the module name in go.mod
go mod edit -module github.com/yourorg/my-new-app

# Download dependencies
go mod tidy
```

Remember to replace `popapp` with your app name all through the codebase, including in `config.yaml`, `mage.go`, `Dockerfile`, and `.gitignore`.

To upgrade dependencies later, bump the Invopop modules explicitly rather than running a blanket `go get -u`: it pulls in a `resty.dev/v3` release candidate that `client.go` does not build against. Keep resty on the version the other Invopop apps use:

```bash
go get github.com/invopop/client.go@latest github.com/invopop/popui.go@latest github.com/invopop/configure@latest
go get resty.dev/v3@v3.0.0-beta.6
go mod tidy
```

### 2. Configuration

Key configuration options:
- `invopop.client_id` and `invopop.client_secret`: Your Invopop app credentials
- `nats.url`: NATS server connection string
- `public_base_url`: Your application's public URL

### 3. Development

#### Using Air (Recommended for Development)

```bash
# Start with hot reload
air
```

This will:
- Watch for file changes
- Automatically rebuild the application
- Restart the server
- Generate Templ templates

#### Using Mage

```bash
# Build the application
mage build

# Start the service
mage serve

# Open a shell in the Docker container
mage shell
```

#### Direct Go Commands

```bash
# Generate Templ templates
templ generate

# Build and run
go build . && ./popapp serve
```

## 🔧 Development Guide

### Adding New Features

The template is organized to make common development tasks straightforward:

#### 1. Adding Web Routes

Edit `internal/interfaces/web/web.go`:

```go
func New(ds *domain.Setup, ic *invopop.Client) *Service {
    s := new(Service)
    s.ic = ic
    s.domain = ds
    s.echo = echopop.NewService()

    s.echo.Serve(func(e *echo.Echo) {
        e.StaticFS(popui.AssetPath, popui.Assets)
        // Add your controllers here. Routes called from the Invopop console
        // should use the authEnrollment middleware to load the enrollment.
        s.app = newAppController(e.Group("/app", s.authEnrollment()), s)
        e.StaticFS("/", assets.Content)
    })

    return s
}
```

#### 2. Adding Gateway Tasks

Edit `internal/interfaces/gateway/gateway.go`:

```go
func (s *Service) executeAction(ctx context.Context, in *gateway.Task) *gateway.TaskResult {
    switch in.Action {
    case "my_new_action":
        return s.handleMyNewAction(ctx, in)
    default:
        return gateway.TaskKO(errors.New("unknown action"))
    }
}
```

#### 3. Adding Domain Logic

Create new files in `internal/domain/`:

```go
// internal/domain/my_service.go
type MyService struct {
    // dependencies
}

func (s *MyService) DoSomething(ctx context.Context, req *MyRequest) (*MyResponse, error) {
    // Business logic here
}
```

#### 4. Adding Templ Components

Create `.templ` files in `internal/interfaces/web/components/`:

```templ
// components/my_component.templ
package components

templ MyComponent(title string) {
    <div class="my-component">
        <h1>{ title }</h1>
    </div>
}
```

## 🐳 Docker Development

The template includes Docker support for consistent development environments:

```bash
# Build and run with Docker
mage serve

# This runs the equivalent of:
docker run --rm --name popapp \
  --network invopop-local \
  -v $PWD:/src -w /src \
  --label traefik.enable=true \
  --label traefik.http.routers.popapp.rule=Host\`popapp.invopop.dev\` \
  --label traefik.http.routers.popapp.tls=true \
  --expose 8080 \
  golang:1.26.4-alpine \
  /src/popapp serve
```

## ⚙️ GitHub Workflows

This template includes three GitHub Actions workflows to automate testing, linting, and deployment:

### 🧪 Test Workflow (`test.yaml`)

**Triggers**: Push to main, tags (v*), and pull requests

**What it does**:
- Sets up Go environment using the version from `go.mod`
- Configures authentication for private Invopop dependencies
- Downloads all Go modules
- Runs unit tests with race detection: `go test -tags unit -race ./...`
- Builds the application to verify compilation

**Required Secrets**:
- `GO_MOD_USER`: GitHub username for accessing private repositories
- `GO_MOD_PASS`: GitHub token/password for private repository access

### 🔍 Lint Workflow (`lint.yaml`)

**Triggers**: Push to main, tags (v*), and pull requests

**What it does**:
- Sets up Go environment
- Configures authentication for private dependencies
- Runs `golangci-lint` to check code quality and style
- Ensures code follows Go best practices

### 🚀 Release Workflow (`release.yaml`)

**Triggers**: Push to main and any tag

**What it does**:
- On a push to main, bumps the semver tag and pushes it
- Builds the Docker image, passing the private-module credentials as build args
- Tags the image with the version and `latest`, and pushes it to Artifact Registry (`europe-docker.pkg.dev/invopop/apps`)

**Required Secrets**:
- `GCP_PROJECT_ID`: Google Cloud project ID
- `GCP_SA_KEY`: JSON service account key with push access to the registry
- `GO_MOD_USER` / `GO_MOD_PASS`: as above, for the Docker build

## 📚 Additional Resources

- [Invopop Documentation](https://docs.invopop.com/)
- [Echo Framework](https://echo.labstack.com/)
- [Templ Documentation](https://templ.guide/)
- [Mage Build Tool](https://magefile.org/)

## 📄 License

This template is provided under the same license as your Invopop application.
