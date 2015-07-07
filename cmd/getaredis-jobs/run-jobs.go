package main

import (
	"log"
	"os"

	"github.com/MohamedBassem/getaredis"
	"github.com/robfig/cron"
)

func monitorHosts() {
	ctx, err := getaredis.Init("../../config.yml")
	errLogger := log.New(os.Stderr, "MonitorHosts", 0)
	outLogger := log.New(os.Stdout, "MonitorHosts", 0)
	outLogger.Println("Started")
	defer outLogger.Println("Finished")
	if err != nil {
		errLogger.Println("Error :" + err.Error())
		return
	}

	started, deleted, err := getaredis.MonitorHosts(ctx)
	if err != nil {
		errLogger.Println("Error :" + err.Error())
		return
	}

	if started {
		outLogger.Println("A new host has started!")
	}

	if deleted != nil && len(deleted) > 0 {
		outLogger.Printf("Hosts %v have been removed.", deleted)
	}
}

func cleanRedisInstances() {
	ctx, err := getaredis.Init("../../config.yml")
	errLogger := log.New(os.Stderr, "CleanRedisInstances", 0)
	outLogger := log.New(os.Stdout, "CleanRedisInstances", 0)
	outLogger.Println("Started")
	defer outLogger.Println("Finished")
	if err != nil {
		errLogger.Println("Error :" + err.Error())
		return
	}

	cleanedInstances := getaredis.CleanRedisInstances(ctx)
	if cleanedInstances != nil && len(cleanedInstances) > 0 {
		outLogger.Printf("Containers %v have been removed.", cleanedInstances)
	}
}

func main() {

	c := cron.New()
	c.AddFunc("@every 20m", cleanRedisInstances)
	c.AddFunc("@every 5m", monitorHosts)
	c.Start()
	<-make(chan struct{})
}
