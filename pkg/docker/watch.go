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

	registryAddr := getDomain(image)
	logrus.Debugf("reg: %s, image: %s", registryAddr, image)

	registry, exist := c.registries[registryAddr]
	if !exist {
		newRegistry, err := reglib.NewFromConfigFile(registryAddr)
		if err != nil {
			return nil, err
		}
		c.registries[registryAddr] = newRegistry
		registry = newRegistry
		logrus.Debugf("set %s registry client", registryAddr)
	}

	watchC := make(chan string)

	go func() {
		defer close(watchC)

		var library bool
		if registryAddr == "index.docker.io" {
			library = true
		}

		index := strings.Index(image, "/")
		image = image[index+1:]
		var (
			repo = image
			tag  = "latest"
		)
		if strings.Contains(image, ":") {
			repo = strings.Split(image, ":")[0]
			tag = strings.Split(image, ":")[1]
		}

		if library {
			repo = "library/" + repo
		}

		logrus.Debugf("watch image: %s:%s", repo, tag)
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
