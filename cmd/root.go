package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "mapper",
		Short: "A sitemap generator tool",
		Long: `Mapper is a command-line tool for generating XML sitemaps by crawling websites.
It stays within the specified domain and supports various configuration options.

Example:
  mapper generate https://example.com`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mapper.yaml)")
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug mode")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (debug, info, warn, error)")
	rootCmd.PersistentFlags().String("user-agent", "", "custom User-Agent string")

	// Bind flags to viper
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.BindPFlag("log_level", rootCmd.PersistentFlags().Lookup("log-level"))
	viper.BindPFlag("user_agent", rootCmd.PersistentFlags().Lookup("user-agent"))

	// Set default values
	viper.SetDefault("debug", false)
	viper.SetDefault("log_level", "info")
	viper.SetDefault("user_agent", "Mapper/1.0 (+https://github.com/ncecere/mapper)")
	viper.SetDefault("concurrent_requests", 5)
	viper.SetDefault("request_timeout", "10s")
	viper.SetDefault("rate_limit", "1s")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".mapper" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mapper")
	}

	// Read in environment variables that match
	viper.AutomaticEnv()
	viper.SetEnvPrefix("MAPPER")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("debug") {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}

// GetUserAgent returns the configured User-Agent string
func GetUserAgent() string {
	return viper.GetString("user_agent")
}

// GetDebugMode returns whether debug mode is enabled
func GetDebugMode() bool {
	return viper.GetBool("debug")
}

// GetLogLevel returns the configured log level
func GetLogLevel() string {
	return viper.GetString("log_level")
}

// GetConcurrentRequests returns the maximum number of concurrent requests
func GetConcurrentRequests() int {
	return viper.GetInt("concurrent_requests")
}

// GetRequestTimeout returns the request timeout duration
func GetRequestTimeout() string {
	return viper.GetString("request_timeout")
}

// GetRateLimit returns the rate limit duration
func GetRateLimit() string {
	return viper.GetString("rate_limit")
}
