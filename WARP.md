# WARP.md

This file provides guidance to WARP (warp.dev) when working with code in this repository.

## Project Overview

Keydra is a HashiCorp Vault auto-unseal utility designed for Kubernetes environments. It monitors Vault's seal status and automatically unseals it using pre-configured unseal keys, ensuring high availability.

## Architecture

The application follows a clean Go project structure:

- **cmd/keydra**: Main application entry point with CLI argument parsing and graceful shutdown
- **pkg/config**: Configuration loading and validation using YAML
- **pkg/vault**: Vault API client wrapper with TLS support
- **pkg/k8s**: Kubernetes client for pod discovery and readiness checks
- **internal/unsealer**: Core business logic that orchestrates the unsealing process

The unsealer runs in a continuous loop, checking Vault's seal status at configured intervals and attempting to unseal when necessary.

## Common Development Commands

### Building and Testing
```bash
# Build the application
make build

# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run single test
go test -v ./pkg/config -run TestLoadConfig

# Format code
make fmt

# Run linter (requires golangci-lint)
make lint

# Clean build artifacts
make clean
```

### Running Locally
```bash
# Run with example config
make run

# Run with custom config file
go run ./cmd/keydra -config=path/to/config.yaml

# Build and run binary
make build && ./bin/keydra -config=config.yaml.example
```

### Docker Operations
```bash
# Build Docker image
make docker-build

# Build and push to registry
DOCKER_REGISTRY=your-registry make docker-push
```

### Kubernetes Development
```bash
# Deploy to local cluster
kubectl apply -f deployments/k8s/rbac.yaml
kubectl apply -f deployments/k8s/deployment.yaml

# View logs
kubectl logs -f deployment/keydra -n vault

# Port forward for metrics
kubectl port-forward svc/keydra 8080:8080 -n vault
```

## Configuration Management

The application uses YAML configuration with three main sections:
- `vault`: Vault connection details, TLS settings, and unseal keys
- `kubernetes`: K8s client configuration for in-cluster or external access  
- `app`: Application settings like check intervals and log levels

Configuration is loaded at startup and validated before creating client instances. The config package handles defaults and validation.

## Key Dependencies

- **github.com/hashicorp/vault/api**: Official Vault Go client
- **k8s.io/client-go**: Kubernetes Go client
- **gopkg.in/yaml.v2**: YAML parsing for configuration

## Security Considerations

- Unseal keys must be stored as Kubernetes secrets, never in code
- Application runs as non-root user (UID 1000) in containers
- Minimal RBAC permissions (read-only access to pods and services)
- TLS enabled by default for Vault connections
- Read-only root filesystem in container

## Testing Strategy

- Unit tests focus on configuration loading and validation
- Integration tests would require Vault and Kubernetes test environments
- Use table-driven tests for configuration validation scenarios
- Test coverage should include error paths and edge cases

## Commit Message Format

Follow the established rule: start every commit with the issue number, e.g. "#1: implement basic unsealer logic"
