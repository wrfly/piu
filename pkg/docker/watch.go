package docker

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wrfly/reglib"
)

func (c *Cli) WatchImageChange(ctx context.Context, image string) (<-chan string, error) {
	if !strings.Contains(image, ":") {
		image += ":latest"
	}

	registryAddr, repo, tag := getMeta(image)

	// docker api starts with index.docker.io
	// but pull image must start with docker.io
	// f**k docker
	if registryAddr == "docker.io" {
		registryAddr = "index.docker.io"
	}

	logrus.Debugf("reg: %s, image: %s", registryAddr, image)

	c.m.RLock()
	registry, exist := c.registries[registryAddr]
	c.m.RUnlock()
	if !exist {
		newRegistry, err := reglib.NewFromConfigFile(registryAddr)
		if err != nil {
			return nil, err
		}
		c.m.Lock()
		c.registries[registryAddr] = newRegistry
		c.m.Unlock()
		registry = newRegistry
		logrus.Debugf("set %s registry client", registryAddr)
	}

	watchC := make(chan string)

	go func() {
		defer close(watchC)

		logrus.Debugf("watch image: %s:%s", repo, tag)
		img, err := registry.Image(ctx, repo, tag)
		if err != nil {
			logrus.Errorf("initial watch image %s error: %s", image, err)
			return
		}

		ident := imageIdentifier(img)
		logrus.Debugf("image %s %s", image, ident)
		for ctx.Err() == nil {
			img, err := registry.Image(ctx, repo, tag)
			if err != nil {
				if ctx.Err() != context.Canceled {
					logrus.Errorf("watch image %s error: %s", image, err)
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
	x := ""
	for _, hist := range img.History() {
		x += hist.Config.Image
	}
	return fmt.Sprintf("%s %s", img.Size(), hash(x))
}

func hash(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}
