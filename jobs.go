package getaredis

import "time"

func CleanRedisInstaces(ctx *context) {
	var instances = make([]Instance, 0)
	var maxTimeStamp = time.Now().Add(-1 * time.Second * 60 * 60 * time.Duration(ctx.config.MaxInstanceTime))
	ctx.db.Model(&Instance{}).Where("running = 1 AND created_at < ?", maxTimeStamp).Find(&instances)

	for _, instance := range instances {
		ctx.RemoveContainer(instance.ContainerID)
	}

	ctx.db.Model(&Instance{}).Where("running = 1 AND created_at < ?", maxTimeStamp).UpdateColumn("running", false)
}
