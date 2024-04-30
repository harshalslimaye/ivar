package main

import (
	"fmt"
	"time"

	"github.com/harshalslimaye/ivar/cli"
	"github.com/harshalslimaye/ivar/cli/initcmd"
	"github.com/harshalslimaye/ivar/cli/installcmd"
	"github.com/logrusorgru/aurora"
)

func main() {
	fmt.Println("⚡️ivar (v0.0.1)")
	timer := time.Now()

	cli.RootCmd.AddCommand(initcmd.InitCmd())
	cli.RootCmd.AddCommand(installcmd.InstallCmd(&timer))

	if err := cli.RootCmd.Execute(); err != nil {
		fmt.Println(aurora.Red("Something went wrong!"))
	}
}
