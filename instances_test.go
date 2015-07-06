package getaredis

import (
	"math/rand"
	"testing"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
)

func getAMockDockerContext() *context {
	rand.Seed(time.Now().UnixNano())
	tmp, _ := docker.NewClient("unix:///var/run/docker.sock")
	ctx := &context{
		dockerClient: *tmp,
	}
	return ctx
}

// TODO Add redis authentication check
func TestStartRedisInstance(t *testing.T) {
	ctx := getAMockDockerContext()
	containerName := generateRandomString(20)
	password := generateRandomString(20)
	container, err := ctx.startRedisInstance(containerName, password)
	if !assert.NoError(t, err, "Starting docker container should not return an Error.") {
		return
	}
	time.Sleep(time.Second)
	container, err = ctx.dockerClient.InspectContainer(containerName)
	if !assert.True(t, container.State.Running, "Container Failed to start.") {
		return
	}
	assert.NotEmpty(t, container.NetworkSettings.Ports["6379/tcp"], "Should have a port mapping for redis port")

	ctx.dockerClient.RemoveContainer(docker.RemoveContainerOptions{
		ID:    container.ID,
		Force: true,
	})
}
