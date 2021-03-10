package task

import (
	v1 "k8s.io/api/core/v1"
	"strings"
)

func getContainerNameIpMapping(pod *v1.Pod)map[string]string{
	cStats := map[string]string{}
	for _,containersSts := range pod.Status.ContainerStatuses{
		containerId := strings.Split(containersSts.ContainerID, "//")
		if len(containerId) != 2 {
			cStats[containersSts.Name] = containersSts.ContainerID
		} else {
			cStats[containersSts.Name] = containerId[1]
		}
	}
	return cStats
}