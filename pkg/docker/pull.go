package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/sirupsen/logrus"
)

func (c *Cli) PullImage(ctx context.Context, image string) error {
	if c.ctx == nil {
		return context.Canceled
	}
	registryAddr, repo, tag := getMeta(image)
	image = fmt.Sprintf("%s/%s:%s", registryAddr, repo, tag)

	logrus.Infof("pulling image %s", image)
	rc, err := c.cli.ImagePull(ctx, image, types.ImagePullOptions{
		RegistryAuth: getAuth(image),
	})
	if err != nil {
		return err
	}
	defer rc.Close()

	ioReader := bufio.NewReader(rc)
	for {
		line, err := ioReader.ReadString('\n')
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}
		logrus.Debugf("pull image %s: %s", image, line)
	}
}
