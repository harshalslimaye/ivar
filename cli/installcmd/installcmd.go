package installcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/Masterminds/semver/v3"
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
			var downloadList sync.Map
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
				downloadList.Store(filepath.Join("node_modules", d.Name()), d)
			}
			WalkGraph(gh, &downloadList)

			var wg sync.WaitGroup
			lr := locker.NewLocker()
			downloadList.Range(func(key, value interface{}) bool {
				downloadPath := key.(string)
				node := value.(*graph.Node)
				lr.Add(
					NewLockItem(node, downloadPath),
					fmt.Sprintf("%s@%s", node.Name(), node.Package.RawVersion),
				)

				wg.Add(1)

				go func(ne *graph.Node, dp string) {
					defer wg.Done()
					if err := DownloadDependency(ne, dp); err != nil {
						fmt.Println(aurora.Red(err))
					} else {
						createSymbolicLink(ne, dp)
					}
				}(node, downloadPath)

				return true
			})

			wg.Wait()
			if err := lr.Write(); err != nil {
				fmt.Println(aurora.Red(err))
			}
			loader.Clear()

			fmt.Printf("%s %s %s\n", "🔥", aurora.Green("success"), "Installation complete!")
			duration := time.Since(t).Round(time.Millisecond * 10)
			fmt.Printf("%s %s %s\n", "⌛", aurora.Cyan("info"), "Done in "+duration.String())
		},
	}

	return cmd
}

func WalkGraph(gh *graph.Graph, dl *sync.Map) {
	defer loader.Clear()

	for _, node := range gh.Nodes {
		WalkNode(nil, node, dl)
	}

}

func WalkNode(parent *graph.Node, node *graph.Node, dl *sync.Map) {
	node.Lock()
	GetDownloadPath(node, parent, dl)
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

func GetDownloadPath(n *graph.Node, p *graph.Node, dl *sync.Map) {
	rootPath := filepath.Join("node_modules", n.Name())
	value, exists := dl.Load(rootPath)

	if exists {
		if node, okay := value.(*graph.Node); !okay || node.Version() != n.Version() {
			dl.Store(filepath.Join("node_modules", p.Name(), "node_modules", n.Name()), n)
			return
		}
	}

	if n.Graph.Versions.Len(n.Name()) > 1 {
		var versions []*semver.Version
		for _, v := range n.Graph.Versions.Get(n.Name()) {
			version, err := semver.NewVersion(v)
			if err != nil {
				fmt.Printf("Error parsing version %s: %v\n", v, err)
				continue
			}
			versions = append(versions, version)
		}

		sort.Slice(versions, func(i, j int) bool {
			return versions[i].LessThan(versions[j])
		})

		if len(versions) > 0 {
			latestVersion := versions[len(versions)-1]
			if latestVersion.String() == n.Version() {
				dl.Store(rootPath, n)
			} else {
				dl.Store(filepath.Join("node_modules", p.Name(), "node_modules", n.Name()), n)
			}
			return
		}
	}

	if !exists {
		dl.Store(rootPath, n)
		return
	}
}

func NewLockItem(node *graph.Node, path string) *locker.Element {
	return &locker.Element{
		Version:   node.Version(),
		Resolved:  node.TarballUrl,
		Integrity: node.Integrity,
		Path:      path,
	}
}
