package getaredis

import "time"

type Address struct {
	IP   string
	port int
}

type Instance struct {
	ID           int
	Name         string
	IP           string
	CreatedAt    time.Time
	HostedAtIP   string
	HostedAtPort int
}

func (ctx *context) NewInstance() Instance {
	instance := *new(Instance)
	return instance
}
func (ctx *context) ListIntances() []Instance {
	instanceList := make([]Instance, 0)
	return instanceList
}
