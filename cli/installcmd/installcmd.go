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
	"github.com/harshalslimaye/ivar/internal/locker"
	"github.com/harshalslimaye/ivar/internal/tarball"
	"github.com/logrusorgru/aurora"
	"github.com/spf13/cobra"
)

func InstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i"},
		Short:   "Installs a package along with its dependencies.",
		Run: func(cmd *cobra.Command, args []string) {
			var downloadList graph.List
			t := time.Now()
			fmt.Println(helper.ShowInfo("📄", "Reading package.json"))
			parser, err := ReadPackageJson()
			if err != nil {
				fmt.Println(aurora.Red(parser))
				os.Exit(1)
			}

			fmt.Println(helper.ShowInfo("🔄", "Resolving Dependencies"))
			gh := graph.NewDependencyGraph(parser)

			if !helper.HasHomeDir() {
				if homedir := helper.HomeDir(); homedir != "" {
					if err := os.MkdirAll(homedir, 0755); err == nil {
						gh.HasCache = true
					}
				}
			} else {
				gh.HasCache = true
			}

			fmt.Println(helper.ShowInfo("📦", "Downloading packages"))
			for _, d := range gh.RootDependencies {
				downloadList.AddNode(d, nil)
			}
			WalkGraph(gh, &downloadList)

			var wg sync.WaitGroup
			Download(downloadList.Map(), &wg)
			wg.Wait()

			fmt.Printf("%s %s %s\n", "🔥", aurora.Green("success"), "Installation complete!")
			duration := time.Since(t).Round(time.Millisecond * 10)
			fmt.Printf("%s %s %s\n", "⌛", aurora.Cyan("info"), "Done in "+duration.String())
		},
	}

	return cmd
}

func Download(nodes map[string]*graph.Node, wg *sync.WaitGroup) {
	for path, node := range nodes {
		wg.Add(1)
		go func(p string, n *graph.Node, wtgp *sync.WaitGroup) {
			defer wtgp.Done()
			if err := DownloadDependency(n, p); err == nil {
				createSymbolicLink(n, p)
				if len(n.PrunedNodes) > 0 {
					Download(n.PrunedNodes, wtgp)
				}
			}
		}(path, node, wg)
	}
}

func WalkGraph(gh *graph.Graph, dl *graph.List) {
	defer loader.Clear()

	for _, node := range gh.Nodes {
		WalkNode(nil, node, dl)
	}

}

func WalkNode(parent *graph.Node, node *graph.Node, dl *graph.List) {
	node.Lock()
	dl.AddNode(node, parent)
	node.Unlock()

	if len(node.Dependencies) > 0 {
		for _, dependencyNode := range node.Dependencies {
			WalkNode(node, dependencyNode, dl)
		}
	}

}
func DownloadDependency(node *graph.Node, dir string) error {
	loader.Show("\r" + "Installing " + node.Name() + "@" + node.Version() + "...")
	if err := tarball.DownloadAndExtract(node, dir); err != nil {
		return err
	}

	return nil
}

func createSymbolicLink(node *graph.Node, dir string) {
	if len(node.Bin) > 0 {
		for name, path := range node.Bin {
			source := filepath.Join(dir, path)
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

func NewLockItem(node *graph.Node, path string) *locker.Element {
	return &locker.Element{
		Version:   node.Version(),
		Resolved:  node.TarballUrl,
		Integrity: node.Integrity,
		Path:      path,
	}
}
