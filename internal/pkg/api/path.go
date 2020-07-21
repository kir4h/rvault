package api

import (
	"fmt"
	"path"
)

// GetListBasePath returns Vault URL base path for listing depending on kv version
func GetListBasePath(engine string, kvVersion string) (string, error) {
	switch kvVersion {
	case "1":
		return engine, nil
	case "2":
		return path.Join(engine, "metadata"), nil
	default:
		return "", fmt.Errorf("KV version '%s' is not supported", kvVersion)
	}
}

// GetReadBasePath returns Vault URL base path for reading depending on kv version
func GetReadBasePath(engine string, kvVersion string) (string, error) {
	switch kvVersion {
	case "1":
		return engine, nil
	case "2":
		return path.Join(engine, "/data"), nil
	default:
		return "", fmt.Errorf("KV version '%s' is not supported", kvVersion)
	}
}
