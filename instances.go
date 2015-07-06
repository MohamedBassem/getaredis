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

func (ctx *context) generateRedisConfig(name, password string) docker.CreateContainerOptions {
	return docker.CreateContainerOptions{
		Name: name,
		Config: &docker.Config{
			Image:      "redis",
			Memory:     5 * 1024 * 1024,
			MemorySwap: -1,
			Cmd:        []string{"redis-server", "--requirepass", password},
		},
	}
}

func (ctx *context) startRedisInstance(name, password string) (*docker.Container, error) {
	container, err := ctx.dockerClient.CreateContainer(ctx.generateRedisConfig(name, password))
	if err != nil {
		return nil, err
	}
	err = ctx.dockerClient.StartContainer(container.ID, &docker.HostConfig{PublishAllPorts: true})
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	container, err = ctx.dockerClient.InspectContainer(container.ID)
	if err != nil || !container.State.Running {
		if err != nil {
			return nil, err
		}
		return nil, errors.New("Container Failed to start")
	}
	return container, nil
}

// Creates a new docker instance with a random name, and returns the instance details back
func (ctx *context) NewInstance(creatorIP, creatorHash string) (*Instance, error) {
	name := generateRandomString(20)
	password := generateRandomString(20)
	var count int
	for ctx.db.Where("name = ?", name).Count(&count); count != 0; name = generateRandomString(20) {
		// Keep Trying!
	}
	container, err := ctx.startRedisInstance(name, password)
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
