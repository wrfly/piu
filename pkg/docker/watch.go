package docker

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/reglib"
)

func (c *Cli) watchImageChange(ctx context.Context, image string) (<-chan string, error) {
	if !strings.Contains(image, ":") {
		image += ":latest"
	}

	registryAddr := getDomain(image)
	logrus.Debugf("reg: %s, image: %s", registryAddr, image)

	c.m.Lock()
	defer c.m.Unlock()
	registry, exist := c.registries[registryAddr]
	if !exist {
		newRegst, err := reglib.NewFromConfigFile(registryAddr)
		if err != nil {
			return nil, err
		}
		c.registries[registryAddr] = newRegst
		registry = newRegst
	}

	watchC := make(chan string)

	go func() {
		defer close(watchC)

		imageURL, _ := url.Parse(image)
		repo := strings.Split(imageURL.Path, ":")[0]
		tag := strings.Split(imageURL.Path, ":")[1]
		img, err := registry.Image(ctx, repo, tag)
		if err != nil {
			logrus.Errorf("watch image error: %s", err)
			return
		}

		ident := imageIdentifier(img)
		logrus.Debugf("image %s %s", image, ident)
		for ctx.Err() == nil {
			img, err := registry.Image(ctx, repo, tag)
			if err != nil {
				if ctx.Err() != context.Canceled {
					logrus.Errorf("watch image error: %s", err)
				}
				return
			}
			newIdent := imageIdentifier(img)
			if ident != newIdent {
				logrus.Infof("image %s changed", image)
				logrus.Debugf("image %s %s", image, newIdent)
				ident = newIdent
				watchC <- newIdent
			}
			time.Sleep(time.Second)
		}
	}()

	return watchC, nil
}

func imageIdentifier(img *reglib.Image) string {
	return fmt.Sprintf("%s %s", img.Size(),
		img.Created().Format(time.RFC3339))
}
