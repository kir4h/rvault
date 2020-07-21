package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"k8s.io/klog/v2"
)

func createDirIfMissing(fs afero.Fs, dir string) error {
	secretDirExists, err := afero.IsDir(fs, dir)
	// path exists and is a folder
	if err == nil && secretDirExists {
		return nil
		// path exists and is a file
	} else if err == nil {
		return fmt.Errorf("unable to create folder %s as it currently exists and is a file", dir)
	}

	// path doesn't exist, create it
	folderPermission := viper.GetUint32("read.folder_permission")
	klog.V(3).Infof("Creating folder '%s'", dir)
	return fs.MkdirAll(dir, os.FileMode(folderPermission))
}

func writeFileContent(fs afero.Fs, filePath string, content string) error {
	filePermission := viper.GetUint32("read.file_permission")
	fo, err := fs.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.FileMode(filePermission))
	if err != nil {
		return err
	}
	_, err = fo.WriteString(content)
	if err != nil {
		return err
	}
	err = fo.Close()
	if err != nil {
		return err
	}
	return nil
}

func dumpToFile(secrets map[string]map[string]string, outputPath string, fs afero.Fs, overwrite bool) error {
	var secretsWritten, secretsSkipped int
	for secret, kv := range secrets {
		secretDir := filepath.Join(outputPath, secret)
		err := createDirIfMissing(fs, secretDir)
		if err != nil {
			return err
		}
		for k, v := range kv {
			filePath := filepath.Join(secretDir, k)
			if !overwrite {
				exists, err := afero.Exists(fs, filePath)
				if err != nil {
					return err
				}
				if exists {
					secretsSkipped += 1
					klog.V(3).Infof("File %s already exists and overwrite is set to 'false', skipping it", filePath)
					continue
				}
			}
			err := writeFileContent(fs, filePath, v)
			if err != nil {
				return err
			}
			secretsWritten += 1
		}
	}
	if secretsWritten != 0 {
		klog.V(2).Infof("Secrets written: %d", secretsWritten)
	}
	if secretsSkipped != 0 {
		klog.V(2).Infof("Secrets skipped: %d", secretsSkipped)
	}

	return nil
}

// Dump writes secrets to the specified format
func Dump(secrets map[string]map[string]string, fs afero.Fs, format string) (result string, err error) {
	switch format {
	case "file":
		outputPath := viper.GetString("read.output")
		overwrite := viper.GetBool("read.overwrite")
		err = dumpToFile(secrets, outputPath, fs, overwrite)
	case "yaml":
		yamlString, err := yaml.Marshal(secrets)
		if err == nil {
			result = string(yamlString)
		}
	case "json":
		jsonString, err := json.Marshal(secrets)
		if err == nil {
			result = string(jsonString)
		}
	default:
		err = fmt.Errorf("unsupported output type '%s'", format)
	}
	return
}
