package getaredis

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/oauth2"

	"github.com/digitalocean/godo"
	"github.com/garyburd/redigo/redis"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

type Host struct {
	IP                 string
	Name               string
	NumberOfContainers int
	MemoryFree         float32
}

func (ctx *context) ListHosts() []Host {
	redisServerKeys, _ := redis.Strings(ctx.redis.Do("KEYS", "server:*"))
	servers := make([]interface{}, len(redisServerKeys))
	for i, t := range redisServerKeys {
		servers[i] = t
	}
	serverConfigs, _ := redis.Strings(ctx.redis.Do("MGET", servers...))

	hosts := make([]Host, len(serverConfigs))
	for i, val := range serverConfigs {
		newHost := new(Host)
		fmt.Println(val)
		err := json.Unmarshal([]byte(val), newHost)
		fmt.Println(err)
		fmt.Println(newHost)
		hosts[i] = *newHost
	}
	fmt.Printf("%+v\n", hosts)
	return hosts
}

func (ctx *context) NewHost() error {
	redisIP := strings.Split(ctx.config.RedisAddress, ":")[0]
	redisPort := strings.Split(ctx.config.RedisAddress, ":")[1]
	dropletName := "getaredis-" + generateRandomString(10)
	userData := `#cloud-config
runcmd:
  - apt-get install -y wget
  - wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
  - tar -C /usr/local -xzf go1.4.2.linux-amd64.tar.gz
  - echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc
  - mkdir /root/go
  - export HOME=/root
  - echo 'export GOROOT=$HOME/go' >> /root/.bashrc
  - echo 'export PATH=$PATH:$GOROOT/bin' >> /root/.bashrc
  - export GOPATH=/root/go
  - /usr/local/go/bin/go get github.com/MohamedBassem/getaredis/...
  - apt-get install -y supervisor
write_files:
  - path: /etc/supervisor/conf.d/go_jobs.conf
    content: |
        [program:go_jobs]
        command=/usr/local/bin/go
        autostart=true
        autorestart=true
        stderr_logfile=/var/log/go_jobs.err.log
        stdout_logfile=/var/log/go_jobs.out.log

        [program:service_discovery]
        command=/usr/local/bin/service_discovery
        autostart=true
        autorestart=true
        stderr_logfile=/var/log/service_discovery.err.log
        stdout_logfile=/var/log/service_discovery.out.log
  - path: /usr/local/bin/service_discovery
    permissions: '0755'
    content: |
        #!/bin/bash
        (
          echo "AUTH %v";
          echo "SET server:%v {}";
          echo "EXPIRE server:%v 10";
          while true; do
            echo "EXPIRE server:%v 10";
            sleep 2;
          done
        ) | telnet %v %v
`

	userData = fmt.Sprintf(userData, ctx.config.RedisPassword, dropletName, dropletName, dropletName, redisIP, redisPort)

	var sshKey *godo.DropletCreateSSHKey
	if ctx.config.DropletSSHKeyID != -1 {
		sshKey = &godo.DropletCreateSSHKey{ID: ctx.config.DropletSSHKeyID}
	}

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: "nyc3",
		Size:   "512mb",
		Image: godo.DropletCreateImage{
			ID: 12380137, // The Docker Image
		},
		UserData: userData,
		SSHKeys:  []godo.DropletCreateSSHKey{*sshKey},
	}

	_, _, err := ctx.digitalocean.Droplets.Create(createRequest)
	return err
}
