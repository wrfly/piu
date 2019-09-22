package docker

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/docker/docker/api/types"
)

func (c *Cli) PullImage(ctx context.Context, image string) error {

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
			if err == io.EOF {
				return nil
			}
			return err
		}
		fmt.Print(line)
	}
}
