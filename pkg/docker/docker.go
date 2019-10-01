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

// Cli with docker
type Cli struct {
	cli client.APIClient
	ctx context.Context

	f map[string]string // list filter
	m sync.RWMutex

	registries map[string]reglib.Registry

	containerChan chan ContainerSpec
}

// New docker cli
func New(ctx context.Context) (*Cli, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}

	cctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	p, err := cli.Ping(cctx)
	if err != nil {
		return nil, err
	}
	logrus.Infof("connect to docker %s", p.APIVersion)

	x := &Cli{
		ctx:           ctx,
		cli:           cli,
		registries:    make(map[string]reglib.Registry),
		containerChan: make(chan ContainerSpec, 100),
	}

	go func() {
		<-ctx.Done()
		cli.Close()
		close(x.containerChan)
	}()

	return x, err
}

// ListContainers and put into container channel
func (c *Cli) ListContainers() error {
	if c.ctx == nil {
		return context.Canceled
	}

	args := filters.NewArgs()
	for k, v := range c.f {
		args.Add(k, v)
	}
	cs, err := c.cli.ContainerList(c.ctx, types.ContainerListOptions{
		Filters: args,
	})
	if err != nil {
		return err
	}

	c.m.Lock()
	for _, container := range cs {
		// NOTE: container can be restarted only got this label
		if container.Labels["piu"] == "" {
			continue
		}
		select {
		case c.containerChan <- ContainerSpec{
			ID:     container.ID,
			Image:  container.Image,
			Action: Start,
		}:
		default:
		}
	}
	c.m.Unlock()

	return nil
}

// Containers channel
func (c *Cli) Containers() chan ContainerSpec {
	return c.containerChan
}
