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

	ctx, cancel := context.WithCancel(context.Background())
	errChan := make(chan error)

	cli, err := docker.New(ctx)
	if err != nil {
		cancel()
		logrus.Fatal("create docker client error: %s", err)
	}

	containerCancelFunc := make(map[string]context.CancelFunc)

	// replace all new create containers with new image
	go func() {
		for container := range cli.Containers() {
			// cancel the old watch goroutine
			if container.Action == docker.Die {
				if cancel, ok := containerCancelFunc[container.ID]; ok {
					logrus.Infof("container %s stopped, stop watching image change", container.ID)
					cancel()
				}
				return
			}
			go func(c docker.ContainerSpec) {
				cCtx, cancel := context.WithCancel(ctx)
				defer cancel()
				containerCancelFunc[c.ID] = cancel

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
				logrus.Infof("stop watching image %s", c.Image)
			}(container)
		}
	}()

	// watch container creation
	go func() {
		if err := cli.WatchEvents(); err != nil {
			errChan <- err
			return
		}
	}()

	// list all the containers
	if err := cli.ListContainers(); err != nil {
		cancel()
		logrus.Fatalf("list containers error: %s", err)
	}

	return watchSig(cancel, errChan)
}

func watchSig(cancel context.CancelFunc, errC chan error) error {
	defer cancel()

	sigC := make(chan os.Signal)
	signal.Notify(sigC, os.Interrupt)

	select {
	case err := <-errC:
		logrus.Fatal("runtime error: %s", err)
		return err
	case s := <-sigC:
		logrus.Warnf("sig %s received, exit", s)
		return nil
	}
}
