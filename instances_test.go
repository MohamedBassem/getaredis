package getaredis

import (
	"testing"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
)

func getDockerClient() (*docker.Client, error) {
	return docker.NewClient("unix:///var/run/docker.sock")
}

func forceRemoveContainer(ctx *context, id string) {
	dockerClient, _ := getDockerClient()
	dockerClient.RemoveContainer(docker.RemoveContainerOptions{
		ID:    id,
		Force: true,
	})
}

// TODO Add redis authentication check
func TestStartRedisInstance(t *testing.T) {
	ctx, _ := Init("config.yml")
	containerName := generateRandomString(20)
	password := generateRandomString(20)
	container, err := startRedisInstance(ctx, "unix:///var/run/docker.sock", containerName, password)
	if !assert.NoError(t, err, "Starting docker container should not return an Error.") {
		return
	}
	time.Sleep(time.Second)
	dockerClient, _ := getDockerClient()
	container, err = dockerClient.InspectContainer(containerName)
	if !assert.True(t, container.State.Running, "Container Failed to start.") {
		return
	}
	assert.NotEmpty(t, container.NetworkSettings.Ports["6379/tcp"], "Should have a port mapping for redis port")
	forceRemoveContainer(ctx, container.ID)
}

// TODO Mock a database for testing and actually test this function
func TestNewInstance(t *testing.T) {
	ctx, _ := Init("config.yml")
	creatorIP, creatorHash := "192.168.1.20", "asdasdgsdasdbdfg"
	instance, _ := ctx.NewInstance(creatorIP)
	forceRemoveContainer(ctx, instance.ContainerID)
}
