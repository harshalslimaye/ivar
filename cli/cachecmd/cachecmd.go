package cachecmd

import (
	"fmt"
	"os"

	"github.com/harshalslimaye/ivar/internal/cache"
	"github.com/harshalslimaye/ivar/internal/filesystem"
	"github.com/harshalslimaye/ivar/internal/helper"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func CacheCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: "Handles packages stored in cache",
		Run: func(cmd *cobra.Command, args []string) {
			if len(args) == 0 {
				fmt.Println(aurora.Red("Please specify the action that needs to be performed"))
				os.Exit(1)
			}

			action := args[0]

			switch action {
			case "clean":
				if err := Clean(); err != nil {
					fmt.Println(aurora.Red(err))
				} else {
					fmt.Println(helper.ShowInfo("🧹", "Cache clean-up complete"))
				}
				break
			default:
				fmt.Println(aurora.Red(fmt.Errorf("invalid action: %s", action)))
				os.Exit(1)
			}
		},
	}

	return cmd
}

func Clean() error {
	if helper.HasHomeDir() {
		c := cache.NewCache()

		if c.IsEmpty() {
			fmt.Println(aurora.Cyan("no packages found in the cache"))
			os.Exit(1)
		}

		if err := filesystem.DeleteContents(helper.HomeDir()); err != nil {
			return err
		}
	}

	return nil
}
