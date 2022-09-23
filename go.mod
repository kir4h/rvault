module rvault

go 1.14

require (
	cloud.google.com/go/kms v1.4.0 // indirect
	cloud.google.com/go/monitoring v1.4.0 // indirect
	github.com/gobwas/glob v0.2.3
	github.com/hashicorp/vault v1.9.9
	github.com/hashicorp/vault-plugin-secrets-kv v0.10.1
	github.com/hashicorp/vault/api v1.3.1
	github.com/hashicorp/vault/sdk v0.3.1-0.20220721224749-00773967ab3a
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.10.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/klog/v2 v2.40.1
)
