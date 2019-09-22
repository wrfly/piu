package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
	"gopkg.in/urfave/cli.v2"

	"github.com/wrfly/piu/pkg/docker"
)

func run(c *cli.Context) error {
	if c.Bool("version") {
		fmt.Println(versionInfo)
		return nil
	}
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	cli, err := docker.New()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)

	// replace all new create containers with new image
	go func() {
		// FIXME: containerChan not closed, goroutine leak
		for container := range cli.Containers() {
			go func(c docker.ContainerSpec) {
				cCtx, cancel := context.WithCancel(ctx)
				defer cancel()

				// FIXME: if container stopped, cancel this goroutine
				changed, err := cli.WatchImageChange(cCtx, c.Image)
				if err != nil {
					logrus.Errorf("watch image [%s] change error: %s", c.Image, err)
					return
				}
				for change := range changed {
					logrus.Infof("container %s changed: %s", c.ID, change)

					// pull the latest image
					if err := cli.PullImage(cCtx, c.Image); err != nil {
						logrus.Warnf("pull image %s error: %s", c.Image, err)
						continue
					}

					// re create the container
					if err := cli.ReCreate(ctx, c.ID); err != nil {
						logrus.Warnf("recreate container %s err: %s", c.ID, err)
					}
				}
			}(container)
		}
	}()

	// watch container creation
	go func() {
		if err := cli.WatchStartEvents(ctx); err != nil {
			errChan <- err
			return
		}
	}()

	// list all the containers
	if err := cli.ListContainers(ctx); err != nil {
		cancel()
		return err
	}

	return watchSig(cancel, errChan)
}

func watchSig(cancel context.CancelFunc, errC chan error) error {
	defer cancel()

	sigC := make(chan os.Signal)
	signal.Notify(sigC, os.Interrupt)

	select {
	case err := <-errC:
		return err
	case <-sigC:
		return nil
	}
}
