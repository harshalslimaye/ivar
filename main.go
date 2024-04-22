package main

import (
	"fmt"

	"github.com/harshalslimaye/ivar/cli"
	"github.com/harshalslimaye/ivar/cli/initcmd"
	"github.com/harshalslimaye/ivar/cli/installcmd"
	"github.com/logrusorgru/aurora"
)

func main() {
	cli.RootCmd.AddCommand(initcmd.InitCmd())
	cli.RootCmd.AddCommand(installcmd.InstallCmd())

	if err := cli.RootCmd.Execute(); err != nil {
		fmt.Println(aurora.Red("Something went wrong!"))
	}
}
