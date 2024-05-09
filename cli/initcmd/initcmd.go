package initcmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/harshalslimaye/ivar/internal/helper"
	"github.com/harshalslimaye/ivar/internal/packagejson"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

const initCmdAlreadyExists = `
Oops! Found a package.json file here already. 
If you want to start fresh, just delete it and 
run 'ivar init' again.
`

const initCmdWelcome = `
%s init walks you through creating a %s file with the essentials, suggesting handy defaults. 
For more info, type 'ivar help init'. Remember, you can press ^C anytime to quit. 
Afterward, use 'ivar install <pkg>' to add dependencies. Easy peasy!
`

func InitCmd() *cobra.Command {
	var yesFlag bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Creates a new package.json file",
		Run: func(cmd *cobra.Command, args []string) {
			pkgjson := packagejson.GetNewPackageJson(yesFlag)

			if packagejson.Exists() {
				fmt.Println(aurora.Red(initCmdAlreadyExists))
				os.Exit(1)
			}

			if !yesFlag {
				fmt.Printf(initCmdWelcome, aurora.Cyan("ivar"), aurora.Green("package.json"))

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
