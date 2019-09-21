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

// Cli with docker
type Cli struct {
	cli client.APIClient
	f   map[string]string

	registries map[string]reglib.Registry
	containers map[string]bool
	m          sync.Mutex
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
		cli:        cli,
		registries: make(map[string]reglib.Registry, 5),
		containers: make(map[string]bool),
	}, err
}

func (c *Cli) listContainers(ctx context.Context) error {
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
		c.containers[container.ID] = true
	}
	c.m.Unlock()

	return nil
}
