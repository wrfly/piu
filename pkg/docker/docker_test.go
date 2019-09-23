package docker

import (
	"context"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func TestDocker(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	cli, err := New(ctx)
	if err != nil {
		t.Fatal(err)
	}

	if err := cli.ListContainers(); err != nil {
		t.Fatal(err)
	}

	for container := range cli.containerChan {
		logrus.Infof("found container %s @%s", container.ID, container.Image)
	}
}

func TestIncreaseContainerName(t *testing.T) {
	tMap := map[string]string{
		"a":           "a.piu.1",
		"hello":       "hello.piu.1",
		"piu":         "piu.piu.1",
		"piu.1":       "piu.1.piu.1",
		"hello.piu.1": "hello.piu.2",
	}
	for k, v := range tMap {

		tv := increaseContainerName(k)
		if tv != v {
			t.Errorf("%s: %s != %s", k, tv, v)
		}
	}
}
