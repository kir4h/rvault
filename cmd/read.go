package cmd

import (
	"fmt"

	"rvault/internal/pkg/api"
	"rvault/internal/pkg/kv"
	"rvault/internal/pkg/output"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"k8s.io/klog/v2"
)

func getFormat(format string) (string, error) {
	switch format {
	case "f", "file":
		return "file", nil
	case "y", "yaml":
		return "yaml", nil
	case "j", "json":
		return "json", nil
	default:
		return "", fmt.Errorf("unsupported format '%s'. Supported formats: 'file', 'yaml', 'json'", format)
	}
}

// readCmd represents the readSecrets command
var readCmd = &cobra.Command{
	Use:   "read <engine>",
	Short: "Recursively read secrets for a given path",
	Long:  "Recursively read secrets for a given path",
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("read.path", cmd.Flags().Lookup("path"))
		_ = viper.BindPFlag("read.output", cmd.Flags().Lookup("output"))
		_ = viper.BindPFlag("read.format", cmd.Flags().Lookup("format"))
		_ = viper.BindPFlag("global.kv_version", cmd.Flags().Lookup("kv-version"))
		_ = viper.BindPFlag("read.overwrite", cmd.Flags().Lookup("overwrite"))
		_ = viper.BindPFlag("read.folder_permission", cmd.Flags().Lookup("folder-permission"))
		_ = viper.BindPFlag("read.file_permission", cmd.Flags().Lookup("file-permission"))

	},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var readSecrets map[string]map[string]string
		c, err := api.NewClient()
		if err != nil {
			klog.Fatalf("%v", err)
		}
		engine := args[0]
		path := viper.GetString("read.path")
		concurrency := viper.GetUint32("global.concurrency")
		includePaths := viper.GetStringSlice("global.include_paths")
		excludePaths := viper.GetStringSlice("global.exclude_paths")
		format, err := getFormat(viper.GetString("read.format"))
		if err != nil {
			klog.Fatalf("%v", err)
		}
		readSecrets, err = kv.RRead(c, engine, path, includePaths, excludePaths, concurrency)
		if err != nil {
			if readSecrets == nil {
				klog.Fatalf("%v", err)
			} else {
				klog.Warningf("%v", err)
			}
		}
		fs := afero.OsFs{}
		res, err := output.Dump(readSecrets, fs, format)
		if err != nil {
			klog.Fatalf("%v", err)
		}
		if res != "" {
			fmt.Println(res)
		}

	},
}

func init() {
	readCmd.Flags().StringP("path", "p", "/", "Path to look for secrets")
	readCmd.Flags().StringP("output", "o", ".", "Output folder for 'file' format ")
	readCmd.Flags().StringP("format", "f", "file", "Output format ('file', 'yaml', 'json')")
	readCmd.Flags().StringP("kv-version", "k", "", "KV Version")
	readCmd.Flags().BoolP("overwrite", "w", false, "Overwrite existing files (file format only)")
	readCmd.Flags().Uint32("folder-permission", 0700, "Permissions for newly created folders (file format only)")
	readCmd.Flags().Uint32("file-permission", 0600, "Permissions for created secret files (file format only)")
	rootCmd.AddCommand(readCmd)
}
