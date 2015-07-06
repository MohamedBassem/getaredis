package getaredis

import (
	"errors"

	"github.com/fsouza/go-dockerclient"
)

import "time"

type Instance struct {
	ID           int `sql:"AUTO_INCREMENT"`
	Name         string
	CreatorIP    string
	CreatorHash  string
	CreatedAt    time.Time
	HostedAtIP   string
	HostedAtPort string
}

func (ctx *context) generateDefaultDockerConfig() *docker.Config {
	return &docker.Config{
		Image:      "redis",
		Memory:     5 * 1024 * 1024,
		MemorySwap: -1,
		PortSpecs:  []string{"6379:6379"},
	}
}

func (ctx *context) generateDockerCreateConfig(name string) docker.CreateContainerOptions {
	return docker.CreateContainerOptions{
		Name:   name,
		Config: ctx.generateDefaultDockerConfig(),
	}
}

func (ctx *context) startDockerInstance(name string) (*docker.Container, error) {
	container, err := ctx.dockerClient.CreateContainer(ctx.generateDockerCreateConfig(name))
	if err != nil {
		return nil, err
	}
	err = ctx.dockerClient.StartContainer(container.ID, &docker.HostConfig{PublishAllPorts: true})
	if err != nil {
		return nil, err
	}
	return ctx.dockerClient.InspectContainer(container.ID)
}

func (ctx *context) NewInstance(creatorIP, creatorHash string) (*Instance, error) {
	name := generateRandomString(20)
	var count int
	for ctx.db.Where("name = ?", name).Count(&count); count != 0; name = generateRandomString(20) {
		// Keep Trying!
	}
	container, err := ctx.startDockerInstance(name)
	if err != nil {
		return nil, err
	}
	instance := &Instance{
		Name:         name,
		CreatorIP:    creatorIP,
		CreatorHash:  creatorHash,
		CreatedAt:    time.Now(),
		HostedAtIP:   container.NetworkSettings.IPAddress,
		HostedAtPort: container.NetworkSettings.Ports["6379/tcp"][0].HostPort,
	}
	ctx.db.Create(instance)
	if ctx.db.NewRecord(instance) {
		return nil, errors.New("Failed to write to the database")
	}
	return instance, nil
}
func (ctx *context) ListIntances() []Instance {
	instanceList := make([]Instance, 0)
	return instanceList
}
