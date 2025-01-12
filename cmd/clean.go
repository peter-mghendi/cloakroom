package cmd

import (
	"cloakroom/lib/handlers"
	"cloakroom/lib/utility"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// cleanCmd represents the clean command
var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Remove all contents from the wardrobe directory.",
	Long: `The clean command removes all files and directories in the specified wardrobe directory.
It is used to prepare for a fresh environment without altering the manifest or lock file.

The wardrobe directory must be defined before running this command.

Example:
  cloakroom clean`,
	Run: func(cmd *cobra.Command, args []string) {
		wardrobe := viper.GetString(utility.Wardrobe)
		err := handlers.Clean(wardrobe)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}
