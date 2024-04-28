package installcmd

import (
	"os"
	"path/filepath"

	"github.com/harshalslimaye/ivar/internal/graph"
	"github.com/harshalslimaye/ivar/internal/helper"
	"github.com/harshalslimaye/ivar/internal/packagejson"
	"github.com/harshalslimaye/ivar/internal/tarball"
	"github.com/spf13/cobra"
)

func InstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "This command installs a package along with its dependencies.",
		Run: func(cmd *cobra.Command, args []string) {
			pkgjson := packagejson.ReadPackageJson()
			gh := graph.NewDependencyGraph(pkgjson.Dependencies)

			WalkGraph(gh)
		},
	}

	return cmd
}

func WalkGraph(gh *graph.Graph) {
	for _, node := range gh.Nodes {
		WalkNode(nil, node, nil)
	}
}

func WalkNode(parent *graph.Node, node *graph.Node, visited map[string]string) {
	if visited == nil {
		visited = make(map[string]string) // keeps track of packages in root node_modules
	}

	// Check if the package has already been visited
	if _, exists := visited[node.Package.Name]; exists {
		if visited[node.Package.Name] != node.Package.Version {
			dir := filepath.Join("node_modules", parent.Package.Name, "node_modules", node.Package.Name)

			// Process the package here (e.g., download and install)
			DownloadDependency(node.Package.Name, node.Package.Version, dir)
		}
	} else {
		dir := filepath.Join("node_modules", node.Package.Name)
		// Mark the package as visited
		visited[node.Package.Name] = node.Package.Version

		// Process the package here (e.g., download and install)
		DownloadDependency(node.Package.Name, node.Package.Version, dir)

		// Recursively walk through dependencies
		for _, dependencyNode := range node.Dependencies {
			WalkNode(node, dependencyNode, visited)
		}
	}
}
func DownloadDependency(name string, version string, downloadDir string) error {
	sourcePath := filepath.Join(downloadDir, name+"-"+version+".tgz")
	targetPath := helper.GetCurrentDirPath() + helper.GetPathSeparator() + downloadDir

	if err := os.MkdirAll(downloadDir, 0755); err != nil {
		return err
	}

	if err := tarball.DownloadTarball(name, version, downloadDir); err != nil {
		return err
	}

	if err := tarball.ExtractTarball(sourcePath, targetPath); err != nil {
		return err
	}

	if err := tarball.DeleteTarball(sourcePath); err != nil {
		return err
	}

	return nil
}
