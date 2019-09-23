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
		domain = "docker.io"
	}
	return domain
}

func getMeta(image string) (domain, repo, tag string) {
	domain = getDomain(image)
	if domain == "docker.io" {
		if !strings.Contains(image, "/") {
			image = "library/" + image
		}
	}
	image = strings.TrimPrefix(image, domain+"/")
	tag = "latest"
	if strings.Contains(image, ":") {
		repo = strings.Split(image, ":")[0]
		tag = strings.Split(image, ":")[1]
	}
	return
}
