package docker

import (
	"context"
	"io"
	"testing"
	"time"
)

func TestPull(t *testing.T) {
	cli, err := New(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = cli.PullImage(ctx, "docker.io/wrfly/piu:echo")
	if err != nil && err != io.EOF {
		t.Errorf("err: %s", err)
	}
}
