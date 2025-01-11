package cmd

import (
	"cloakroom/lib"
	"cloakroom/lib/handlers"
	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add <owner/repo>",
	Short: "Add a new plugin to the manifest.",
	Long: `The add command adds a new plugin to the manifest.

You can specify the plugin's repository and release details. Optionally, use the --fetch flag to immediately download the plugin after adding it.

Example:
  cloakroom add example/my-plugin --release v1.2.0 --artifact plugin.jar
  cloakroom add example/my-plugin --release latest --artifact plugin.jar --fetch`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		wardrobe := viper.GetString("WARDROBE")
		manifest := &lib.Manifest{}
		err := viper.Unmarshal(manifest)
		cobra.CheckErr(err)

		key := args[0]
		release, _ := cmd.Flags().GetString("release")
		artifact, _ := cmd.Flags().GetString("artifact")
		fetch, _ := cmd.Flags().GetBool("fetch")

		plugin := lib.Plugin{
			Tag:      release,
			Artifact: artifact,
		}

		err = handlers.Add(manifest, plugin, key, wardrobe, fetch)
		cobra.CheckErr(err)

		manifest.Plugins[key] = plugin
		viper.Set("plugins", manifest.Plugins)
		err = viper.WriteConfig()
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	addCmd.Flags().String("release", "", "Tag version of the plugin (required).")
	addCmd.Flags().String("artifact", "", "Artifact name of the plugin (required).")
	addCmd.Flags().Bool("fetch", false, "Immediately download the plugin after adding it.")
}
