module rvault

go 1.14

require (
	github.com/gobwas/glob v0.2.3
	github.com/hashicorp/vault v1.4.3
	github.com/hashicorp/vault-plugin-secrets-kv v0.5.5
	github.com/hashicorp/vault/api v1.0.5-0.20200317185738-82f498082f02
	github.com/hashicorp/vault/sdk v0.1.14-0.20200702114606-96dd7d6e10db
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/afero v1.6.0
	github.com/spf13/cobra v1.1.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/klog v1.0.0 // indirect
	k8s.io/klog/v2 v2.40.1
)
