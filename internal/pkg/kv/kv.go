package kv

import (
	"fmt"
	"strings"

	vapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

func parseSecretData(result readResult, kvVersion string) (secretData map[string]string, errString string) {
	var data map[string]interface{}
	var ok bool

	secretData = make(map[string]string)

	dataValue, err := getDataValueFromResult(result, kvVersion)
	if err != nil {
		errString = err.Error()
		return
	}

	if dataValue == nil {
		klog.V(5).Infof("Discarding empty secret %s", result.path)
		return
	}

	if data, ok = dataValue.(map[string]interface{}); !ok {
		errString = fmt.Sprintf("Discarding secret %s as its data field is of unexpected type %T", result.path,
			dataValue)
	}

	for k, v := range data {
		secretData[k] = fmt.Sprintf("%v", v)
	}

	return
}

func getDataValueFromResult(result readResult, kvVersion string) (interface{}, error) {
	if kvVersion == "1" {
		return result.secret.Data, nil
	}

	if dataValue, ok := result.secret.Data["data"]; !ok {
		return nil, fmt.Errorf("discarding secret '%s' because 'data' field was not found", result.path)
	} else {
		return dataValue, nil
	}
}

func getMountOutput(c *vapi.Client, engine string) (*vapi.MountOutput, error) {
	mounts, err := c.Sys().ListMounts()
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(engine, "/") {
		engine = engine + "/"
	}
	if mount, ok := mounts[engine]; !ok {
		return nil, fmt.Errorf("engine %s not found", engine)
	} else {
		return mount, nil
	}
}

func getKVVersion(c *vapi.Client, engine string) (string, error) {
	// look for engine specific setting
	version := viper.GetString(fmt.Sprintf("engines.%s.kv_version", engine))
	if version == "" {
		// fallback to global value
		version = viper.GetString("global.kv_version")
	}

	if version != "" {
		klog.V(4).Infof("Using kv version '%s' for engine '%s'", version, engine)
		return version, nil
	}
	klog.V(3).Infof("KV Version not specified, listing mounts to find it")
	// Version not specified, listing mounts to get the version
	mc, err := getMountOutput(c, engine)
	if err != nil {
		return "", err
	}

	if mc.Type != "kv" {
		return "", fmt.Errorf("unsupported engine type '%s'", mc.Type)
	}
	return mc.Options["version"], nil
}
