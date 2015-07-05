package getaredis

import "time"

type Instance struct {
	ID           int `sql:"AUTO_INCREMENT"`
	Name         string
	CreatorIP    string
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
