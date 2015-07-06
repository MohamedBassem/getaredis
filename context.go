package getaredis

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/fsouza/go-dockerclient"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type configuration struct {
	Database        map[string]string `yaml:"database"`
	DockerHost      string            `yaml:"dockerHost"`
	MaxInstanceSize int               `yaml:"maxInstanceSize"`
	MaxInstanceTime int               `yaml:"maxInstanceTime"`
}

type context struct {
	dockerClient docker.Client
	db           gorm.DB
	config       configuration
}

func Init(configPath string) (*context, error) {
	rand.Seed(time.Now().UnixNano())
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := new(configuration)
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	tmp, _ := docker.NewClient(config.DockerHost)
	databaseHost := fmt.Sprintf("%v:%v@/%v?charset=utf8&parseTime=True&loc=Local", config.Database["user"], config.Database["password"], config.Database["dbname"])
	tmp2, err := gorm.Open("mysql", databaseHost)
	if err != nil {
		return nil, err
	}
	tmp2.AutoMigrate(&Instance{})
	ctx := context{
		dockerClient: *tmp,
		config:       *config,
		db:           tmp2,
	}
	ctx.db.LogMode(true)
	return &ctx, nil
}
