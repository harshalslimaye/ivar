package cli

import (
	"fmt"

	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "ivar",
	Short: "Ivar - A JavaScript package manager",
	Long:  "Ivar is a simple JavaScript package manager built with Go.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(aurora.Green("Welcome to Ivar!"))
		_ = cmd.Help()
	},
}
