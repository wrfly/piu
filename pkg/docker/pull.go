package docker

import (
	"bufio"
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
)

func (c *Cli) pullImage(ctx context.Context, image string) error {

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
		if err != nil {
			return err
		}
		fmt.Print(line)
	}
}
