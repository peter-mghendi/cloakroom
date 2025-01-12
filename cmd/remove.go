package cmd

import (
	"cloakroom/lib"
	"cloakroom/lib/handlers"
	"cloakroom/lib/utility"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <artifact>",
	Short: "Remove a plugin from the manifest.",
	Long: `The remove command removes a plugin from the manifest.

Optionally, use the --purge flag to delete the plugin's files from the wardrobe directory.

Examples:
  cloakroom remove plugin.jar
  cloakroom remove plugin.jar --purge`,
	Args: cobra.ExactArgs(1), // Requires exactly one argument: artifact name
	Run: func(cmd *cobra.Command, args []string) {
		wardrobe := viper.GetString(utility.Wardrobe)
		manifest := &lib.Manifest{}
		err := viper.Unmarshal(manifest)
		cobra.CheckErr(err)

		key := args[0]
		purge, _ := cmd.Flags().GetBool("purge")

		err = handlers.Remove(manifest, key, wardrobe, purge)
		cobra.CheckErr(err)

		viper.Set("plugins", manifest.Plugins)
		err = viper.WriteConfig()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)

	// Add flags for the remove command
	removeCmd.Flags().Bool("purge", false, "Delete plugin files from the wardrobe after removal.")
}
