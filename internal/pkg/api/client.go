package api

import (
	"fmt"
	"strings"

	vapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

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

// NewClient returns a new Client
func NewClient() (*vapi.Client, error) {
	config, err := getVaultConfig()
	if err != nil {
		return nil, err
	}
	vaultToken, err := getVaultToken()
	if err != nil {
		return nil, err
	}

	c, err := vapi.NewClient(config)
	if err != nil {
		return nil, err
	}
	klog.V(4).Infof("Vault Adddress: '%s'", c.Address())

	c.SetToken(vaultToken)
	return c, nil
}
