package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of CtrApi",
	Long:  `All software has versions. This is CtrApi's`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("CtrApi Static Site Generator v0.1 -- HEAD")
	},
}
