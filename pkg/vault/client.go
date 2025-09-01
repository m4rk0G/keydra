package vault

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/hashicorp/vault/api"
)

// Client wraps the Vault API client
type Client struct {
	client *api.Client
	keys   []string
}

// NewClient creates a new Vault client
func NewClient(address string, keys []string, tlsConfig TLSConfig) (*Client, error) {
	config := api.DefaultConfig()
	config.Address = address

	if tlsConfig.Enabled {
		tlsClientConfig := &tls.Config{
			InsecureSkipVerify: tlsConfig.Insecure,
		}

		if tlsConfig.CACert != "" {
			// TODO: Load CA certificate
		}

		if tlsConfig.ClientCert != "" && tlsConfig.ClientKey != "" {
			// TODO: Load client certificate
		}

		transport := &http.Transport{TLSClientConfig: tlsClientConfig}
		config.HttpClient.Transport = transport
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}

	return &Client{
		client: client,
		keys:   keys,
	}, nil
}

// TLSConfig represents TLS configuration for Vault client
type TLSConfig struct {
	Enabled    bool
	CACert     string
	ClientCert string
	ClientKey  string
	Insecure   bool
}

// IsSealed checks if Vault is sealed
func (c *Client) IsSealed() (bool, error) {
	status, err := c.client.Sys().SealStatus()
	if err != nil {
		return false, fmt.Errorf("failed to get seal status: %w", err)
	}
	return status.Sealed, nil
}

// Unseal attempts to unseal Vault using the configured keys
func (c *Client) Unseal() error {
	for _, key := range c.keys {
		resp, err := c.client.Sys().Unseal(key)
		if err != nil {
			return fmt.Errorf("failed to unseal with key: %w", err)
		}

		if !resp.Sealed {
			return nil // Successfully unsealed
		}
	}

	return fmt.Errorf("failed to unseal Vault with provided keys")
}

// GetSealStatus returns the current seal status
func (c *Client) GetSealStatus() (*api.SealStatusResponse, error) {
	return c.client.Sys().SealStatus()
}
