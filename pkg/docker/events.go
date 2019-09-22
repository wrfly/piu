package docker

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

func (c *Cli) WatchStartEvents(ctx context.Context) error {
	startFilter := filters.NewArgs()
	startFilter.Add("event", "start")
	msgC, errC := c.cli.Events(ctx, types.EventsOptions{
		Filters: startFilter,
	})
	go func() {
		for msg := range msgC {
			containerID := msg.ID
			image := msg.Actor.Attributes["image"]
			select {
			case c.containerChan <- ContainerSpec{
				ID:    containerID,
				Image: image,
			}:
			default:
			}
		}
	}()

	return <-errC
}
