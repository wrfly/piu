package docker

import (
	"bufio"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/wrfly/reglib"
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

func getAuth(image string) string {
	domain := strings.Split(image, "/")[0]
	if !strings.Contains(domain, ".") {
		domain = "index.docker.io"
	}
	user, pass := reglib.GetAuthFromFile(domain)
	auth := types.AuthConfig{
		Username: user,
		Password: pass,
	}
	authBytes, _ := json.Marshal(auth)
	return base64.URLEncoding.EncodeToString(authBytes)
}
