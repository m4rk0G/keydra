package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary config file
	configContent := `
vault:
  nodes:
    - "https://vault1.example.com:8200"
	- "https://vault2.example.com:8200"
  namespace: "example-namespace"
  unseal_keys:
    - "key1"
    - "key2"
  tls:
    enabled: true
    insecure: false

kubernetes:
  in_cluster: true
  namespace: "vault"

app:
  check_interval: "60s"
  log_level: "debug"
`

	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Load the config
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify loaded values
	if cfg.Vault.Nodes[0] != "https://vault1.example.com:8200" {
		t.Errorf("Expected vault address 'https://vault1.example.com:8200', got '%s'", cfg.Vault.Nodes[0])
	}

	if len(cfg.Vault.Keys) != 2 {
		t.Errorf("Expected 2 unseal keys, got %d", len(cfg.Vault.Keys))
	}

	if cfg.App.CheckInterval != 60*time.Second {
		t.Errorf("Expected check interval 60s, got %v", cfg.App.CheckInterval)
	}
}

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				Vault: VaultConfig{
					Nodes: []string{"https://vault.example.com"},
					Keys:  []string{"key1", "key2"},
				},
			},
			wantErr: false,
		},
		{
			name: "missing vault address",
			config: &Config{
				Vault: VaultConfig{
					Keys: []string{"key1", "key2"},
				},
			},
			wantErr: true,
		},
		{
			name: "missing unseal keys",
			config: &Config{
				Vault: VaultConfig{
					Nodes: []string{"https://vault.example.com"},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
