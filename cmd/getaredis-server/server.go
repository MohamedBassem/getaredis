package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/MohamedBassem/getaredis"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
)

func main() {
	configFileName := flag.String("config", "", "Configuration file path")
	port := flag.Int("port", 8080, "Server listening port")
	flag.Parse()

	if *configFileName == "" {
		log.Fatal("A configuration file must be provided.")
	}

	ctx, err := getaredis.Init(*configFileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	m := martini.Classic()
	m.Use(render.Renderer())

	m.Get("/", func(r render.Render) {
		r.HTML(200, "index", "")
	})

	m.Post("/instance", func(res http.ResponseWriter, req *http.Request) (int, string) {
		requesterIP := strings.Split(req.RemoteAddr, ":")[0]
		err := getaredis.CheckInstanceLimit(ctx, requesterIP)
		if err != nil {
			return 403, err.Error()
		}
		instance, err := ctx.NewInstance(requesterIP)
		if err != nil {
			return 500, err.Error()
		}
		return 200, fmt.Sprintf("{\"IP\": \"%v\", \"port\": \"%v\", \"password\": \"%v\"}", instance.HostedAtIP, instance.HostedAtPort, instance.Password)
	})
	m.RunOnAddr("127.0.0.1:" + strconv.Itoa(*port))
}
