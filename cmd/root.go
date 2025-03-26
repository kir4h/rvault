package cmd

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/klog/v2"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var (
	version = "DEV"
	commit  = "unknown"
	date    = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "rvault",
	Version: fmt.Sprintf("%s \nCommit: %s\nDate: %s", version, commit, date),
	Short:   "Tool to perform some recursive operations on Vault KV",
	Long:    `Tool to perform some recursive operations on Vault KV`,
}

// Execute runs Cobra's root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	klog.InitFlags(nil)

	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("v"))
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("logtostderr"))
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("alsologtostderr"))
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("log_file"))
	pflag.CommandLine.AddGoFlag(flag.CommandLine.Lookup("log_dir"))

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "",
		"config file (default is $HOME/.config/rvault/config.yaml)")
	rootCmd.PersistentFlags().StringP("address", "a", "", "Vault address")
	rootCmd.PersistentFlags().StringP("token", "t", "", "Vault token")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "Vault namespace")
  rootCmd.PersistentFlags().Uint32P("concurrency", "c", 20, "Maximum number of concurrent queries to Vault")
	rootCmd.PersistentFlags().StringSliceP("include-paths", "i", []string{"*"}, "KV paths to be included")
	rootCmd.PersistentFlags().StringSliceP("exclude-paths", "e", []string{},
		"KV paths to be excluded (applied on 'include-paths' output)")
	rootCmd.PersistentFlags().Bool("insecure", false, "Disables SSL verification")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".config/rvault/config.yaml".
		viper.AddConfigPath(home + "/.config/rvault")
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	_ = viper.BindEnv("global.address", "VAULT_ADDR")
	_ = viper.BindEnv("global.token", "VAULT_TOKEN")
	_ = viper.BindEnv("global.namespace", "VAULT_NAMESPACE")
	_ = viper.BindPFlag("global.address", rootCmd.Flags().Lookup("address"))
	_ = viper.BindPFlag("global.token", rootCmd.Flags().Lookup("token"))
	_ = viper.BindPFlag("global.namespace", rootCmd.Flags().Lookup("namespace"))
	_ = viper.BindPFlag("global.verbosity", rootCmd.Flags().Lookup("v"))
	_ = viper.BindPFlag("global.concurrency", rootCmd.Flags().Lookup("concurrency"))
	_ = viper.BindPFlag("global.include_paths", rootCmd.Flags().Lookup("include-paths"))
	_ = viper.BindPFlag("global.exclude_paths", rootCmd.Flags().Lookup("exclude-paths"))
	_ = viper.BindPFlag("global.insecure", rootCmd.Flags().Lookup("insecure"))
	_ = flag.Set("v", viper.GetString("global.verbosity"))
	if err == nil {
		klog.V(1).Infof("Using config file: '%s'", viper.ConfigFileUsed())
	}

	// Read ~/.vault-token file
	if !viper.IsSet("global.token") {
		home, err := homedir.Dir()
		if err != nil {
			return
		}
		tokenFile := home + "/.vault-token"
		if _, err = os.Stat(tokenFile); !os.IsNotExist(err) {
			fh, err := os.Open(tokenFile)
			if err != nil {
				klog.Exitf("Can't read token from file '%s': %v", tokenFile, err)
			}
			defer fh.Close()

			s := bufio.NewScanner(fh)
			if s.Scan() {
				viper.SetDefault("global.token", s.Text())
			}
		}
	}
}
