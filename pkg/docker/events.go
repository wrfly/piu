package docker

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/sirupsen/logrus"
)

// WatchEvents (start and die)
func (c *Cli) WatchEvents() error {
	startFilter := filters.NewArgs()
	startFilter.Add("event", "start")
	startFilter.Add("event", "die")
	msgC, errC := c.cli.Events(c.ctx, types.EventsOptions{
		Filters: startFilter,
	})
	go func() {
		for msg := range msgC {
			containerID := msg.ID
			image := msg.Actor.Attributes["image"]
			select {
			case c.containerChan <- ContainerSpec{
				ID:     containerID,
				Image:  image,
				Action: Action(msg.Action),
			}:
			default:
				logrus.Warnf("container events %v missed", msg)
			}
		}
	}()

	return <-errC
}
