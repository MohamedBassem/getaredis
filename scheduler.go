package getaredis

import "errors"

// Returns the ip for the host which can hold this container
func (ctx *context) scheduleNewContainer() (publicIP, privateIP string, err error) {
	hosts := ctx.ListHosts()
	maximumNumberOfContainers := ctx.config.MaxContainersPerHost
	maximumNumber := -1
	chosenHost := -1
	for i, host := range hosts {
		if host.NumberOfContainers < maximumNumberOfContainers && maximumNumber < host.NumberOfContainers {
			maximumNumber = host.NumberOfContainers
			chosenHost = i
		}
	}

	if chosenHost == -1 {
		return "", "", errors.New("Cannot schedule container.")
	}
	return hosts[chosenHost].PublicIP, hosts[chosenHost].PrivateIP, nil
}
