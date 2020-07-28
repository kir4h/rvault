package kv_test

import (
	"testing"

	kv "github.com/hashicorp/vault-plugin-secrets-kv"
	vapi "github.com/hashicorp/vault/api"
	"github.com/hashicorp/vault/http"
	"github.com/hashicorp/vault/sdk/logical"
	"github.com/hashicorp/vault/vault"
)

var secrets map[string]map[string]interface{}
var engine = "secret"
var engineV2 = "secretv2"

func init() {
	secrets = map[string]map[string]interface{}{
		"/spain/admin": {
			"admin.conf": "dsfdsflfrf43l4tlp",
		},
		"/spain/malaga/random": {
			"my.key": "d3ewf2323r21e2",
		},
		"/france/paris/key": {
			"id_rsa": "ewdfpelfr23pwlrp32l4[p23lp2k",
			"id_dsa": "fewfowefkfkwepfkewkfpweokfeowkfpk",
		},
		"/france/nantes/dummy": {
			"": "",
		},
		"/uk/london/mi5": {
			"mi5.conf": "salt, 324r23432, false",
		},
	}
}

func writeSecret(t *testing.T, c *vapi.Client, path string, content map[string]interface{}) {
	t.Helper()
	_, err := c.Logical().Write(path, content)
	if err != nil {
		t.Fatal(err)
	}
}

func deleteSecret(t *testing.T, c *vapi.Client, path string) {
	t.Helper()
	_, err := c.Logical().Delete(path)
	if err != nil {
		t.Fatal(err)
	}
}

func createSecretsV1(t *testing.T, c *vapi.Client, engine string, secrets map[string]map[string]interface{}) {
	t.Helper()
	for path, content := range secrets {
		writeSecret(t, c, engine+path, content)
	}
}

func createSecretsV2(t *testing.T, c *vapi.Client, engine string, secrets map[string]map[string]interface{}) {
	t.Helper()
	for path, content := range secrets {
		writeSecret(t, c, engine+"/data"+path, map[string]interface{}{"data": content})
	}
}

func createTestVault(t *testing.T) *vault.TestCluster {
	t.Helper()

	coreConfig := &vault.CoreConfig{
		LogicalBackends: map[string]logical.Factory{
			"kv": kv.Factory,
		},
	}
	cluster := vault.NewTestCluster(t, coreConfig, &vault.TestClusterOptions{
		HandlerFunc: http.Handler,
	})
	cluster.Start()
	client := cluster.Cores[0].Client
	err := cluster.Cores[0].Client.Sys().Mount("secretv2", &vapi.MountInput{
		Type:    "kv",
		Options: map[string]string{"version": "2"},
	})
	if err != nil {
		t.Fatal(err)
	}

	createSecretsV1(t, client, "secret", secrets)
	deleteSecret(t, client, "secret/france/nantes/dummy")
	createSecretsV2(t, client, "secretv2", secrets)
	deleteSecret(t, client, "secretv2/metadata/france/nantes/dummy")

	return cluster
}
