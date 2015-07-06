package getaredis

import (
	"math/rand"
	"time"

	"github.com/fsouza/go-dockerclient"
	"github.com/jinzhu/gorm"
)

type context struct {
	dockerClient docker.Client
	db           gorm.DB
}

func Init() context {
	rand.Seed(time.Now().UnixNano())
	tmp, _ := docker.NewClient("unix:///var/run/docker.sock")
	ctx := context{
		dockerClient: *tmp,
	}
	return ctx
}
