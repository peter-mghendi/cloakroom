package cmd

import (
	"cloakroom/lib"
	"cloakroom/lib/handlers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all plugins currently defined in the manifest.",
	Long: `The list command outputs all the plugin definitions present in your Cloakroom manifest.

It prints essential information such as the release tag, artifact name, and an optional hash for verification. If no plugins are found, it displays a simple message.`,
	Run: func(cmd *cobra.Command, args []string) {
		manifest := &lib.Manifest{}
		err := viper.Unmarshal(manifest)
		cobra.CheckErr(err)

		err = handlers.List(manifest)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
