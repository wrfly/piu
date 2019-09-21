package docker

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestWatchImageChange(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)

	cli, err := New()
	if err != nil {
		t.Fatal(err)
	}

	image := "wrfly/hello-world"

	changed, err := cli.watchImageChange(context.Background(), image)
	if err != nil {
		t.Fatal(err)
	}
	for x := range changed {
		t.Logf("image %s changed to %s", image, x)
	}
}
