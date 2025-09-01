package unsealer

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/m4rk0G/keydra/pkg/config"
	"github.com/m4rk0G/keydra/pkg/k8s"
	"github.com/m4rk0G/keydra/pkg/vault"
)

// Unsealer handles the auto-unsealing of Vault
type Unsealer struct {
	vaultClients []*vault.Client
	k8sClient    *k8s.Client
	config       *config.Config
}

// New creates a new Unsealer instance
func New(cfg *config.Config) (*Unsealer, error) {
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Create Vault client
	vaultTLS := vault.TLSConfig{
		Enabled:    cfg.Vault.TLS.Enabled,
		CACert:     cfg.Vault.TLS.CACert,
		ClientCert: cfg.Vault.TLS.ClientCert,
		ClientKey:  cfg.Vault.TLS.ClientKey,
		Insecure:   cfg.Vault.TLS.Insecure,
	}

	var vaultClients []*vault.Client
	for _, node := range cfg.Vault.Nodes {
		vaultClient, err := vault.NewClient(node, cfg.Vault.Keys, vaultTLS)
		if err != nil {
			return nil, fmt.Errorf("failed to create Vault client: %w", err)
		}
		vaultClients = append(vaultClients, vaultClient)
	}

	// Create Kubernetes client
	k8sClient, err := k8s.NewClient(cfg.K8s.InCluster, cfg.K8s.ConfigPath, cfg.K8s.Namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return &Unsealer{
		vaultClients: vaultClients,
		k8sClient:    k8sClient,
		config:       cfg,
	}, nil
}

// Start begins the auto-unsealing loop
func (u *Unsealer) Start(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Printf("Starting auto-unsealer with check interval: %v", u.config.App.CheckInterval)

	ticker := time.NewTicker(u.config.App.CheckInterval)
	defer ticker.Stop()

	// Initial check
	if err := u.checkAndUnseal(ctx); err != nil {
		log.Printf("Initial unseal check failed: %v", err)
	}

	// Periodic checks
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutdown signal received, stopping auto-unsealer")
			return nil
		case <-ticker.C:
			if err := u.checkAndUnseal(ctx); err != nil {
				log.Printf("Unseal check failed: %v", err)
			}
		}
	}
}

// checkAndUnseal performs the actual check and unsealing logic
func (u *Unsealer) checkAndUnseal(ctx context.Context) error {
	// Check if Vault pods are ready
	ready, err := u.k8sClient.IsVaultReady(ctx)
	if err != nil {
		return fmt.Errorf("failed to check Vault readiness: %w", err)
	}

	if !ready {
		log.Println("No Vault pods are ready, skipping unseal check")
		return nil
	}

	// Check if any Vault is sealed
	for _, vaultClient := range u.vaultClients {
		sealed, err := vaultClient.IsSealed()
		if err != nil {
			return fmt.Errorf("failed to check seal status: %w", err)
		}

		if !sealed {
			log.Println("Vault is already unsealed")
			return nil
		}

		log.Println("Vault is sealed, attempting to unseal...")

		// Attempt to unseal
		if err := vaultClient.Unseal(); err != nil {
			return fmt.Errorf("failed to unseal Vault: %w", err)
		}

		log.Println("Successfully unsealed Vault")
	}

	return nil
}
