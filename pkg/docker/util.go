package docker

import (
	"encoding/base64"
	"encoding/json"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/wrfly/reglib"
)

func getAuth(image string) string {
	domain := getDomain(image)
	user, pass := reglib.GetAuthFromFile(domain)
	auth := types.AuthConfig{
		Username: user,
		Password: pass,
	}
	authBytes, _ := json.Marshal(auth)
	return base64.URLEncoding.EncodeToString(authBytes)
}

func getDomain(image string) string {
	domain := strings.Split(image, "/")[0]
	if !strings.Contains(domain, ".") {
		domain = "index.docker.io"
	}
	return domain
}
