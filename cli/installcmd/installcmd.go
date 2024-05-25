package installcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/harshalslimaye/ivar/internal/cmdshim"
	"github.com/harshalslimaye/ivar/internal/graph"
	"github.com/harshalslimaye/ivar/internal/helper"
	"github.com/harshalslimaye/ivar/internal/jsonparser"
	"github.com/harshalslimaye/ivar/internal/loader"
	"github.com/harshalslimaye/ivar/internal/tarball"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func InstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "This command installs a package along with its dependencies.",
		Run: func(cmd *cobra.Command, args []string) {
			t := time.Now()
			fmt.Println(helper.ShowInfo("📄", "Reading package.json"))
			parser, err := ReadPackageJson()
			if err != nil {
				fmt.Println(aurora.Red(parser))
				os.Exit(1)
			}

			fmt.Println(helper.ShowInfo("🔄", "Resolving Dependencies"))
			gh := graph.NewDependencyGraph(parser)

			fmt.Println(helper.ShowInfo("📦", "Fetching packages"))
			WalkGraph(gh)
			fmt.Printf("%s %s %s\n", "🔥", aurora.Green("success"), "Installation complete!")
			duration := time.Since(t).Round(time.Millisecond * 10)
			fmt.Printf("%s %s %s\n", "⌛", aurora.Cyan("info"), "Done in "+duration.String())
		},
	}

	return cmd
}

func WalkGraph(gh *graph.Graph) {
	var visited sync.Map
	var wg sync.WaitGroup
	defer loader.Clear()

	for _, node := range gh.Nodes {
		WalkNode(nil, node, &visited, &wg)
	}

	wg.Wait()
}

func WalkNode(parent *graph.Node, node *graph.Node, visited *sync.Map, wg *sync.WaitGroup) {
	loader.Show("\r" + "Installing " + node.Name() + "@" + node.Version() + "...")

	wg.Add(1)
	go func() {
		defer wg.Done()

		version, exists := visited.Load(node.Package.Name)

		if exists {
			if version != node.Version() {
				node.DownloadDir = filepath.Join("node_modules", parent.Name(), "node_modules", node.Name())
				if helper.Exists(node.DownloadDir) {
					node.DownloadDir = ""
				}
			}
		} else {
			node.DownloadDir = filepath.Join("node_modules", node.Name())
			visited.Store(node.Package.Name, node.Version())
			if helper.Exists(node.DownloadDir) {
				if helper.SameVersionExists(node.DownloadDir, node.Version()) {
					node.DownloadDir = ""
				} else {
					alternateDir := filepath.Join("node_modules", parent.Name(), "node_modules", node.Name())
					if helper.Exists(alternateDir) {
						node.DownloadDir = ""
					} else {
						node.DownloadDir = alternateDir
					}
				}
			}
		}

		if node.DownloadDir != "" {
			if err := DownloadDependency(node); err != nil {
				fmt.Println(aurora.Red(err))
			} else {
				createSymbolicLink(node)
				for _, dependencyNode := range node.Dependencies {
					WalkNode(node, dependencyNode, visited, wg)
				}
			}
		}
	}()
}
func DownloadDependency(node *graph.Node) error {
	if err := tarball.DownloadAndExtract(node); err != nil {
		return err
	}

	return nil
}

func createSymbolicLink(node *graph.Node) {
	if len(node.Bin) > 0 {
		for name, path := range node.Bin {
			source := filepath.Join(node.DownloadDir, path)
			target := filepath.Join("node_modules", ".bin", name)

			cmdshim.CmdShim(source, target)
		}
	}
}

func ReadPackageJson() (*jsonparser.JsonParser, error) {
	data, err := os.ReadFile(helper.GetPackageJsonPath())
	if err != nil {
		return nil, fmt.Errorf("unable to read package.json: %s", err.Error())
	}

	return jsonparser.NewJsonParserFromBytes(data)
}
