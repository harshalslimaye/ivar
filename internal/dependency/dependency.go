package dependency

type Dependency struct {
	Name         string
	Version      string
	Dependencies []Dependency
}

type Tree struct {
	Dependencies []Dependency
}
