package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new configuration file.",
	Long: `The init command creates a new configuration file with default settings.

If a configuration file already exists, it will not overwrite it unless the --force flag is provided.

Example:
  cloakroom init
  cloakroom init --force`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")

		defaultConfig := map[string]interface{}{
			"version": "1.0",
			"host":    "github.com",
		}

		if _, err := os.Stat(viper.ConfigFileUsed()); err == nil {
			if !force {
				fmt.Printf("[ERROR] Config file already exists: %s\n", viper.ConfigFileUsed())
				fmt.Println("Use the --force flag to overwrite the existing configuration.")
				return
			} else {
				fmt.Printf("[INFO] Overwriting config file: %s\n", viper.ConfigFileUsed())
				_ = os.Remove(viper.ConfigFileUsed())
			}
		}

		viper.Reset()
		for key, value := range defaultConfig {
			viper.Set(key, value)
		}

		if err := viper.WriteConfigAs("./cloakroom.json"); err != nil {
			fmt.Printf("[ERROR] Failed to write config file: %v\n", err)
			return
		}

		fmt.Println("[INFO] Configuration file initialized successfully.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolP("force", "f", false, "Overwrite existing configuration file if it exists.")
}
