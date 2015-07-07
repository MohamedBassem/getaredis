package getaredis

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"time"

	"golang.org/x/oauth2"

	"gopkg.in/yaml.v2"

	"github.com/digitalocean/godo"
	"github.com/garyburd/redigo/redis"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type configuration struct {
	Database            map[string]string `yaml:"database"`
	RedisAddress        string            `yaml:"redisAddress"`
	RedisPassword       string            `yaml:"redisPassword"`
	DigitalOceanToken   string            `yaml:"digitalOceanToken"`
	DropletSSHKeyID     int               `yaml:"dropletSSHKeyID"`
	MaxInstanceSize     int               `yaml:"maxInstanceSize"`
	MaxInstanceTime     int               `yaml:"maxInstanceTime"`
	MaxInstancesPerIP   int               `yaml:"maxInstancesPerIP"`
	MaxRedisConnections int               `yaml:"maxRedisConnections"`
}

type context struct {
	db           gorm.DB
	redis        redis.Conn
	digitalocean godo.Client
	config       configuration
}

func Init(configPath string) (*context, error) {
	rand.Seed(time.Now().UnixNano())
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	config := new(configuration)
	config.DropletSSHKeyID = -1 // Default Value
	err = yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	// Starting mysql connection
	databaseHost := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8&parseTime=True&loc=Local", config.Database["user"], config.Database["password"], config.Database["host"], config.Database["dbname"])
	tmp2, err := gorm.Open("mysql", databaseHost)
	if err != nil {
		return nil, err
	}
	tmp2.AutoMigrate(&Instance{})

	// Starting redis connection
	tmp3, err := redis.Dial("tcp", config.RedisAddress)
	if err != nil {
		return nil, err
	}
	tmp3.Do("AUTH", config.RedisPassword)

	// Starting digital ocean client
	oauthClient := oauth2.NewClient(oauth2.NoContext, &TokenSource{AccessToken: config.DigitalOceanToken})
	tmp4 := godo.NewClient(oauthClient)

	ctx := context{
		config:       *config,
		db:           tmp2,
		redis:        tmp3,
		digitalocean: *tmp4,
	}

	return &ctx, nil
}
