package getaredis

import (
	"encoding/json"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

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
