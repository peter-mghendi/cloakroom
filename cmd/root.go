package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloakroom",
	Short: "Manage Keycloak plugins with ease and consistency.",
	Long: `Cloakroom is a CLI tool for managing Keycloak plugins across environments.

It allows you to:
- Initialize a manifest for your plugin configuration.
- Add, remove, or restore plugins based on the manifest and lock file.
- Clean plugin directories and ensure consistency across multiple environments.

Examples:
  # Restore plugins from the manifest
  cloakroom restore

  # Add a new plugin to the manifest
  cloakroom add user/repo --release latest --artifact plugin.jar

  # Clean plugin directory
  cloakroom clean
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(configure)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "Specify a custom config file (default is ./cloakroom.json)")
}

// configure reads in configure file and ENV variables if set.
func configure() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigType("json")
		viper.SetConfigName("cloakroom")
	}

	viper.SetEnvPrefix("CLOAKROOM")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err == nil {
		_, _ = fmt.Fprintf(os.Stdout, "[INFO] Using config file: %s\n", viper.ConfigFileUsed())
		return
	}

	var configFileNotFoundError viper.ConfigFileNotFoundError
	if errors.As(err, &configFileNotFoundError) {
		cmd, _, _ := rootCmd.Find(os.Args[1:])
		if cmd != nil && cmd.Name() == "init" {
			_, _ = fmt.Println("[INFO] No config file found. This is expected for 'init' command.")
			return
		}

		_, _ = fmt.Println("[ERROR] Config file not found. Run 'cloakroom init' to create one.")
		os.Exit(1)
	}

	cobra.CheckErr(err)
}
