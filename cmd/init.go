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
	Short: "Initialize a new manifest file.",
	Long: `The init command creates a new manifest file with default settings.

If a manifest file already exists, it will not overwrite it unless the --force flag is provided.

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
				fmt.Printf("[ERROR] Manifest already exists: %s\n", viper.ConfigFileUsed())
				fmt.Println("Use the --force flag to overwrite the existing manifest.")
				return
			} else {
				fmt.Printf("[INFO] Overwriting manifest: %s\n", viper.ConfigFileUsed())
				_ = os.Remove(viper.ConfigFileUsed())
			}
		}

		viper.Reset()
		for key, value := range defaultConfig {
			viper.Set(key, value)
		}

		if err := viper.WriteConfigAs("./cloakroom.json"); err != nil {
			fmt.Printf("[ERROR] Failed to write manifest: %v\n", err)
			return
		}

		fmt.Println("[INFO] Configuration file initialized successfully.")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().Bool("force", false, "Overwrite existing manifest file if it exists.")
}
