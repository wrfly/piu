package docker

// Config ...
type Config struct {
	Path    string
	Filters map[string]string
}

// Action of the container
type Action string

// actions
const (
	Start Action = "start"
	Die   Action = "die"
)

// ContainerSpec the infomation of the container
type ContainerSpec struct {
	ID     string
	Image  string
	Action Action
}
