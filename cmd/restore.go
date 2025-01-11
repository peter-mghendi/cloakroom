package cmd

import (
	"cloakroom/lib"
	"cloakroom/lib/handlers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var clean bool
var force bool

// restoreCmd represents the restore command
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Sync the local environment with the plugin manifest and lock file.",
	Long: `The restore command ensures that your local environment is synchronized with the defined plugin manifest.

By default, it will:
- Check for an existing lock file and install plugins listed in it.
- Generate a new lock file if none exists, based on the manifest.
- Skip installation of plugins that already exist in the target directory.

Flags:
- Use the --clean (-c) flag to delete existing plugin directories before restoring, ensuring a fresh environment.
- Use the --force (-f) flag to overwrite plugin directories even if they already exist.

Examples:
  # Standard restore
  cloakroom restore

  # Restore with a clean slate
  cloakroom restore --clean

  # Restore and force overwrite of existing plugins
  cloakroom restore --force

  # Clean and force restore
  cloakroom restore --clean --force
`,
	Run: func(cmd *cobra.Command, args []string) {
		manifest := &lib.Manifest{}
		err := viper.Unmarshal(manifest)

		cobra.CheckErr(err)

		wardrobe := viper.GetString("WARDROBE")
		err = handlers.Restore(manifest, wardrobe, clean, force)
		cobra.CheckErr(err)
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)

	restoreCmd.Flags().BoolVarP(&clean, "clean", "c", false, "Remove all plugins before restoring.")
	restoreCmd.Flags().BoolVarP(&force, "force", "f", false, "Overwrite existing plugin directories.")
}
