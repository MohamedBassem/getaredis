package getaredis

import "github.com/fsouza/go-dockerclient"

type context struct {
	dockerClient docker.Client
}
