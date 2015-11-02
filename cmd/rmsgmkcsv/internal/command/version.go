package command

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of " + CommandName + ".",
	Long:  `All software has versions. This is ` + CommandName + `'s`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(CommandName + " v1 -- HEAD")
	},
}
