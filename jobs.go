package getaredis

import "time"

func CleanRedisInstaces(ctx *context) (NumberOfContainerStopped int) {
	var instances = make([]Instance, 0)
	var maxTimeStamp = time.Now().Add(-1 * time.Second * 60 * 60 * time.Duration(ctx.config.MaxInstanceTime))
	ctx.db.Model(&Instance{}).Where("running = 1 AND created_at < ?", maxTimeStamp).Find(&instances)

	for _, instance := range instances {
		ctx.RemoveContainer(instance.HostedAtIP, instance.ContainerID)
	}
	ctx.db.Model(&Instance{}).Where("running = 1 AND created_at < ?", maxTimeStamp).UpdateColumn("running", false)

	return len(instances)
}

func MonitorHosts(ctx *context) (startedHosts, deletedHosts int, err error) {
	hosts := ctx.ListHosts()
	zeros := 0
	for _, host := range hosts {
		if host.NumberOfContainers == 0 {
			zeros++
		}
	}
	if zeros == 0 || len(hosts) == 0 {
		err = ctx.NewHost()
		if err != nil {
			return -1, -1, err
		}
		startedHosts++
		return
	} else if zeros > 1 && len(hosts) > 2 {
		for _, host := range hosts {
			if zeros == 1 {
				break
			}
			if host.NumberOfContainers == 0 {
				err = ctx.DeleteHost(host.IP)
				zeros--
				deletedHosts++
			}
		}
	}
	return
}
