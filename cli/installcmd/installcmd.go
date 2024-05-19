package installcmd

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"

	cmdShim "github.com/harshalslimaye/ivar/internal/cmd-shim"
	"github.com/harshalslimaye/ivar/internal/graph"
	"github.com/harshalslimaye/ivar/internal/helper"
	"github.com/harshalslimaye/ivar/internal/packagejson"
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
			pkgjson := packagejson.ReadPackageJson()

			fmt.Println(helper.ShowInfo("🔄", "Resolving Dependencies"))
			gh := graph.NewDependencyGraph(pkgjson.GetProjectDependencies())

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
	defer fmt.Print("\r\033[K")

	for _, node := range gh.Nodes {
		WalkNode(nil, node, &visited, &wg)
	}

	wg.Wait()
}

func WalkNode(parent *graph.Node, node *graph.Node, visited *sync.Map, wg *sync.WaitGroup) {
	showInstalling(node)

	wg.Add(1)
	// Check if the package has already been visited
	go func() {
		defer wg.Done()

		version, exists := visited.Load(node.Package.Name)

		if exists {
			if version != node.Package.Version {
				node.DownloadDir = filepath.Join("node_modules", parent.Package.Name, "node_modules", node.Package.Name)
			}
		} else {
			node.DownloadDir = filepath.Join("node_modules", node.Package.Name)
			visited.Store(node.Package.Name, node.Package.Version)
		}

		// Process the package here (e.g., download and install)
		if err := DownloadDependency(node); err != nil && node.DownloadDir != "" {
			fmt.Println(aurora.Red(err))
		} else {
			createSymbolicLink(node)
		}

		// Recursively walk through dependencies
		for _, dependencyNode := range node.Dependencies {
			WalkNode(node, dependencyNode, visited, wg)
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

			cmdShim.CmdShim(source, target)
		}
	}
}

func showInstalling(node *graph.Node) {
	fmt.Print("\r\033[K")
	fmt.Print("\r" + "Installing " + node.Name() + "@" + node.Version() + "...")
}
