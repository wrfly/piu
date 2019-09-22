package docker

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/network"
)

// ReCreate a container
func (c *Cli) ReCreate(ctx context.Context, cid string) error {
	// inspect
	cJSON, err := c.cli.ContainerInspect(ctx, cid)
	if err != nil {
		return err
	}
	if !cJSON.State.Running {
		logrus.Warnf("container %s not running, stop re-creating", cid)
		return nil
	}

	// reset configs
	cJSON.Image = ""
	// FIXME: diff image env and runtime env
	// image cmd and runtime cmd ...
	cJSON.Config.Env = nil
	cJSON.Config.Cmd = nil

	resp, err := c.cli.ContainerCreate(ctx,
		cJSON.Config,
		cJSON.HostConfig,
		&network.NetworkingConfig{},
		increaseContainerName(cJSON.Name),
	)
	logrus.Infof("create container %s", resp.ID)

	// stop the old one
	second5 := time.Second * 5
	if err := c.cli.ContainerStop(ctx, cid, &second5); err != nil {
		return fmt.Errorf("stop container %s err: %s", cid, err)
	}

	// start the new one
	return c.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
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
