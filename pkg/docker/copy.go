package docker

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

// Copy a container
func (c *Cli) Copy(ctx context.Context, cid string) error {
	// inspect
	cJSON, err := c.cli.ContainerInspect(ctx, cid)
	if err != nil {
		return err
	}

	resp, err := c.cli.ContainerCreate(ctx,
		cJSON.Config,
		cJSON.HostConfig,
		&network.NetworkingConfig{},
		increaseContainerName(cJSON.Name),
	)

	logrus.Infof("create container %s", resp.ID)

	c.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})

	return err
}

var suffix = regexp.MustCompile(`\.piu\.[0-9]$`)

func increaseContainerName(s string) string {
	if strings.Contains(s, "piu") {
		if suffix.MatchString(s) {
			v, err := strconv.Atoi(string(s[len(s)-1]))
			if err != nil {
				return s + ".piu.1"
			}
			return fmt.Sprintf("%s.%d", s[:len(s)-2], v+1)
		}
	}
	return s + ".piu.1"
}
