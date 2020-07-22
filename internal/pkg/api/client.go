package api

import (
	"fmt"
	"strings"

	vapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"k8s.io/klog"
)

// VaultClient is an interface that describes a vault Client.
type VaultClient interface {
	List(path string) (*vapi.Secret, error)
	Read(path string) (*vapi.Secret, error)
	ListMounts() (map[string]*vapi.MountOutput, error)
}

// Client is an implementation of a vault Client interface.
type Client struct {
	VC *vapi.Client
}

func getVaultConfig() (*vapi.Config, error) {
	vaultAddress := viper.GetString("global.address")
	if vaultAddress == "" {
		return nil, fmt.Errorf("missing mandatory vault address")
	}
	if !strings.HasPrefix(vaultAddress, "http") && !strings.HasPrefix(vaultAddress, "https") {
		vaultAddress = "http://" + vaultAddress
	}
	klog.V(5).Infof("Vault address: %s", vaultAddress)
	config := &vapi.Config{Address: vaultAddress}
	err := config.ConfigureTLS(&vapi.TLSConfig{
		Insecure: viper.GetBool("global.insecure"),
	})
	return config, err
}

func getVaultToken() (string, error) {
	vaultToken := viper.GetString("global.token")
	if vaultToken == "" {
		return "", fmt.Errorf("missing mandatory vault token")
	}
	return vaultToken, nil
}

// List lists Vault secrets
func (c *Client) List(path string) (*vapi.Secret, error) {
	return c.VC.Logical().List(path)
}

// Read reads Vault secrets
func (c *Client) Read(path string) (*vapi.Secret, error) {
	return c.VC.Logical().Read(path)
}

// ListMounts list Vault Mounts
func (c *Client) ListMounts() (map[string]*vapi.MountOutput, error) {
	return c.VC.Sys().ListMounts()
}

// NewClient returns a new Client
func NewClient() (*Client, error) {
	config, err := getVaultConfig()
	if err != nil {
		return nil, err
	}
	vaultToken, err := getVaultToken()
	if err != nil {
		return nil, err
	}

	vc, err := vapi.NewClient(config)
	if err != nil {
		return nil, err
	}

	c := &Client{vc}
	c.VC.SetToken(vaultToken)
	return c, nil
}
