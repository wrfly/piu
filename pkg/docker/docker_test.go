package docker

import (
	"context"
	"testing"
	"time"
)

func TestDocker(t *testing.T) {
	cli, err := New()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	cli.listContainers(ctx)

	for cid := range cli.containers {
		cli.Copy(ctx, cid)
		break
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
