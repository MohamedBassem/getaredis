package getaredis

import (
	"encoding/json"
	"errors"
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
	PublicIP           string
	PrivateIP          string
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
		err := json.Unmarshal([]byte(val), newHost)
		if err != nil {
			continue
		}
		hosts[i] = *newHost
	}
	return hosts
}

func (ctx *context) NewHost() error {
	redisIP := strings.Split(ctx.config.RedisAddress, ":")[0]
	redisPort := strings.Split(ctx.config.RedisAddress, ":")[1]
	dropletName := "getaredis-server-" + generateRandomString(10)
	userData := `#cloud-config
runcmd:
  - docker pull redis
  - apt-get install -y supervisor nginx apache2-utils
  - echo 'DOCKER_OPTS=$DOCKER_OPTS" -H unix:///var/run/docker.sock -H tcp://127.0.0.1:2375"' >> /etc/default/docker
  - service docker restart
  - mkdir -p /etc/nginx/docker_auth
  - htpasswd -b -c /etc/nginx/docker_auth/.htpasswd %v %v
  - service nginx reload
write_files:
  - path: /etc/supervisor/conf.d/go_jobs.conf
    content: |
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
          PUBLIC_IP=$(curl http://169.254.169.254/metadata/v1/interfaces/public/0/ipv4/address)
          PRIVATE_IP=$(curl http://169.254.169.254/metadata/v1/interfaces/private/0/ipv4/address)
          NODE_NAME=%v
          echo "AUTH %v";
          while true; do
            NUMBER_OF_CONTAINERS=$(($(docker ps | wc -l) - 1))
            echo "SET server:$NODE_NAME '{\"PublicIP\":\"$PUBLIC_IP\",\"PrivateIP\":\"$PRIVATE_IP\",\"Name\":\"$NODE_NAME\",\"NumberOfContainers\":$NUMBER_OF_CONTAINERS}'";
            echo "EXPIRE server:$NODE_NAME 10";
            sleep 4;
          done
        ) | telnet %v %v
  - path : /etc/nginx/conf.d/docker.conf
    content: |
        server {
          listen 2377;
          location / {
            auth_basic "Restricted";
            auth_basic_user_file /etc/nginx/docker_auth/.htpasswd;
            proxy_buffering off;
            proxy_pass http://localhost:2375;
          }
        }
`

	userData = fmt.Sprintf(userData, ctx.config.Docker["user"], ctx.config.Docker["password"], dropletName, ctx.config.RedisPassword, redisIP, redisPort)

	var sshKey *godo.DropletCreateSSHKey
	if ctx.config.DropletSSHKeyID != -1 {
		sshKey = &godo.DropletCreateSSHKey{ID: ctx.config.DropletSSHKeyID}
	}

	createRequest := &godo.DropletCreateRequest{
		Name:   dropletName,
		Region: "nyc3",
		Size:   "512mb",
		Image: godo.DropletCreateImage{
			Slug: "docker",
		},
		UserData:          userData,
		PrivateNetworking: true,
		SSHKeys:           []godo.DropletCreateSSHKey{*sshKey},
	}

	_, _, err := ctx.digitalocean.Droplets.Create(createRequest)
	return err
}

func (ctx *context) DeleteHost(ip string) error {
	droplets, _, err := ctx.digitalocean.Droplets.List(nil)
	if err != nil {
		return err
	}
	deleted := false
	for _, d := range droplets {
		for _, intr := range d.Networks.V4 {
			if intr.IPAddress == ip {
				_, err := ctx.digitalocean.Droplets.Delete(d.ID)
				if err != nil {
					return err
				}
				deleted = true
				break
			}
		}
	}
	if !deleted {
		return errors.New("Couldn't find droplet with this IP")
	}
	return nil
}
