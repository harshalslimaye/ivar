package initcmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/harshalslimaye/ivar/internal/helper"
	"github.com/harshalslimaye/ivar/internal/messages"
	"github.com/harshalslimaye/ivar/internal/packagejson"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func InitCmd() *cobra.Command {
	var yesFlag bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates a new package.json file",
		Run: func(cmd *cobra.Command, args []string) {
			pkgjson := packagejson.GetNewPackageJson(yesFlag)

			if packagejson.Exists() {
				fmt.Println(aurora.Red(messages.InitCmdAlreadyExists))
				os.Exit(1)
			}

			if !yesFlag {
				fmt.Printf(messages.InitCmdWelcome, aurora.Cyan("ivar"), aurora.Green("package.json"))

				pkgjson.Name = helper.AskQuestion(`package name: %s`, helper.GetCurrentDirName())
				pkgjson.Version = helper.AskQuestion(`version %s`, "1.0.0")
				pkgjson.Description = helper.AskQuestion(`description`, "")
				pkgjson.Main = helper.AskQuestion("entry point %s", "index.js")
				pkgjson.Repository = helper.AskQuestion("git repository", "")
				pkgjson.License = helper.AskQuestion("license %s", "MIT")

				fmt.Println(aurora.Sprintf(aurora.Cyan("Getting ready to write to %s"), aurora.Green(helper.GetPackageJsonPath())))

				pkgjson.PrintInitJson()

				if okay := helper.AskQuestion("Is this OK? %s", "yes"); strings.ToLower(okay) == "yes" || strings.ToLower(okay) == "y" {
					pkgjson.WriteToFile(helper.GetFileName())
				} else {
					fmt.Println(aurora.Red("Oops, operation canceled!"))
				}
			} else {
				pkgjson.WriteToFile(helper.GetFileName())
			}
		},
	}

	cmd.Flags().BoolVarP(&yesFlag, "yes", "y", false, "Skip interactive prompts and use default values")

	return cmd
}
