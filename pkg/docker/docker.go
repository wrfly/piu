package docker

import (
	"context"
	"sync"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/sirupsen/logrus"
	"github.com/wrfly/reglib"
)

// Config ...
type Config struct {
	Path    string
	Filters map[string]string
}

type ContainerSpec struct {
	ID    string
	Image string
}

// Cli with docker
type Cli struct {
	cli client.APIClient

	f map[string]string // list filter
	m sync.Mutex

	registries map[string]reglib.Registry

	containerChan chan ContainerSpec
}

// New docker cli
func New() (*Cli, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	p, err := cli.Ping(ctx)
	if err != nil {
		return nil, err
	}
	logrus.Infof("connect to %v", p)
	return &Cli{
		cli:           cli,
		registries:    make(map[string]reglib.Registry),
		containerChan: make(chan ContainerSpec, 100),
	}, err
}

func (c *Cli) ListContainers(ctx context.Context) error {
	args := filters.NewArgs()
	for k, v := range c.f {
		args.Add(k, v)
	}
	cs, err := c.cli.ContainerList(ctx, types.ContainerListOptions{
		Filters: args,
	})
	if err != nil {
		return err
	}

	c.m.Lock()
	for _, container := range cs {
		logrus.Infof("found container: %s, image: %s",
			container.ID, container.Image)
		select {
		case c.containerChan <- ContainerSpec{
			ID:    container.ID,
			Image: container.Image,
		}:
		default:
		}
	}
	c.m.Unlock()

	return nil
}

func (c *Cli) Containers() chan ContainerSpec {
	return c.containerChan
}
