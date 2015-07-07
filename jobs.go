package getaredis

import "time"

func CleanRedisInstances(ctx *context) (containerNames []string) {
	var instances = make([]Instance, 0)
	containerNames = make([]string, 0)
	var maxTimeStamp = time.Now().Add(-1 * time.Second * 60 * 60 * time.Duration(ctx.config.MaxInstanceTime))
	ctx.db.Model(&Instance{}).Where("running = 1 AND created_at < ?", maxTimeStamp).Find(&instances)

	for _, instance := range instances {
		ctx.RemoveContainer(instance.HostedAtIP, instance.ContainerID)
		containerNames = append(containerNames, instance.Name)
	}
	ctx.db.Model(&Instance{}).Where("running = 1 AND created_at < ?", maxTimeStamp).UpdateColumn("running", false)

	return
}

func MonitorHosts(ctx *context) (startedHosts bool, deletedHosts []string, err error) {
	hosts := ctx.ListHosts()
	deletedHosts = make([]string, 0)
	zeros := 0
	for _, host := range hosts {
		if host.NumberOfContainers == 0 {
			zeros++
		}
	}
	if zeros == 0 || len(hosts) < 2 {
		err = ctx.NewHost()
		if err != nil {
			return
		}
		if len(hosts) == 0 {
			err = ctx.NewHost()
			if err != nil {
				return
			}
		}
		startedHosts = true
		return
	} else if zeros > 1 && len(hosts) > 2 {
		for _, host := range hosts {
			if zeros == 1 {
				break
			}
			if host.NumberOfContainers == 0 {
				err = ctx.DeleteHost(host.PublicIP)
				zeros--
				deletedHosts = append(deletedHosts, host.PublicIP)
			}
		}
	}
	return
}
