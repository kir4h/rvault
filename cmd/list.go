package cmd

import (
	"fmt"

	"rvault/internal/pkg/api"
	"rvault/internal/pkg/kv"

	"k8s.io/klog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the listSecrets command
var listCmd = &cobra.Command{
	Use:   "list <engine>",
	Short: "Recursively list secrets for a given path",
	Long:  `Recursively list secrets for a given path`,
	PreRun: func(cmd *cobra.Command, args []string) {
		_ = viper.BindPFlag("list.path", cmd.Flags().Lookup("path"))
		_ = viper.BindPFlag("global.kv_version", cmd.Flags().Lookup("kv-version"))
	},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		c, err := api.NewClient()
		if err != nil {
			klog.Fatalf("%v", err)
		}
		engine := args[0]
		path := viper.GetString("list.path")
		concurrency := viper.GetUint32("global.concurrency")
		includePaths := viper.GetStringSlice("global.include_paths")
		excludePaths := viper.GetStringSlice("global.exclude_paths")
		secrets, err := kv.RList(c, engine, path, includePaths, excludePaths, concurrency)
		if err != nil {
			klog.Fatalf("%v", err)
		}
		for _, secret := range secrets {
			fmt.Printf("%s\n", secret)
		}

	},
}

func init() {
	listCmd.Flags().StringP("path", "p", "/", "Path to look for secrets")
	listCmd.Flags().StringP("kv-version", "k", "", "KV Version")
	rootCmd.AddCommand(listCmd)
}
