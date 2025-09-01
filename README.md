# Keydra

Keydra is an auto-unseal utility for HashiCorp Vault running in Kubernetes. It continuously monitors Vault's seal status and automatically unseals it when necessary, ensuring high availability of your Vault cluster.

## Features

- **Automatic Unsealing**: Monitors Vault's seal status and automatically unseals when needed
- **Kubernetes Native**: Designed to run seamlessly in Kubernetes environments
- **Configurable**: Flexible YAML-based configuration
- **Secure**: Follows security best practices with minimal privileges
- **Monitoring**: Built-in metrics endpoint for observability
- **Graceful Shutdown**: Handles termination signals properly

## Quick Start

### Prerequisites

- Go 1.21 or later
- Access to a Kubernetes cluster
- HashiCorp Vault instance running in Kubernetes
- Vault unseal keys

### Building

```bash
# Build the binary
make build

# Run tests
make test

# Build Docker image
make docker-build
```

### Configuration

Copy the example configuration and customize it:

```bash
cp config.yaml.example config.yaml
# Edit config.yaml with your Vault details and unseal keys
```

### Running Locally

```bash
# Run with custom config
go run ./cmd/keydra -config=config.yaml

# Or using make
make run
```

### Deploying to Kubernetes

1. Create the necessary secrets and configmaps:

```bash
# Create namespace
kubectl create namespace vault

# Create secret with unseal keys
kubectl create secret generic vault-unseal-keys \
  --from-literal=key1="your-base64-unseal-key-1" \
  --from-literal=key2="your-base64-unseal-key-2" \
  --from-literal=key3="your-base64-unseal-key-3" \
  -n vault

# Create configmap with configuration
kubectl create configmap keydra-config \
  --from-file=config.yaml \
  -n vault
```

2. Deploy Keydra:

```bash
# Deploy RBAC
kubectl apply -f deployments/k8s/rbac.yaml

# Deploy application
kubectl apply -f deployments/k8s/deployment.yaml
```

## Configuration

Keydra uses YAML configuration. Key sections:

- **vault**: Vault connection settings and unseal keys
- **kubernetes**: Kubernetes client configuration
- **app**: Application-specific settings like check intervals

See `config.yaml.example` for a complete example.

## Security Considerations

- Store unseal keys securely using Kubernetes secrets
- Use RBAC to limit Keydra's permissions
- Enable TLS for Vault connections
- Run with non-root user (default in Docker image)
- Use read-only root filesystem in containers

## Development

### Project Structure

```
keydra/
├── cmd/keydra/           # Main application entry point
├── pkg/
│   ├── config/           # Configuration management
│   ├── vault/            # Vault client wrapper
│   └── k8s/              # Kubernetes client wrapper
├── internal/unsealer/    # Core unsealing logic
├── deployments/k8s/      # Kubernetes deployment manifests
└── Makefile              # Build automation
```

### Available Make Targets

- `make build` - Build the binary
- `make test` - Run tests
- `make test-coverage` - Run tests with coverage
- `make lint` - Run linter
- `make fmt` - Format code
- `make clean` - Clean build artifacts
- `make docker-build` - Build Docker image
- `make help` - Show all available targets

## License

MIT License - see [LICENSE](LICENSE) file for details.
