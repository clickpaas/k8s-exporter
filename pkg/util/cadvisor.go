package util


// doc
// all unit from cadvisor will trans to kb

import (
	"github.com/google/cadvisor/client"
	cadvisorv1 "github.com/google/cadvisor/info/v1"
)

// MachineInfo represent information of an physical node
type MachineInfo struct {
	// total memory of machine
	MemCapacity uint64
}

// ContainerInfo represent running information of container
type ContainerInfo struct {
	// MemRss represent the memory is used by container, not virtualMemory
	MemRss uint64
	MemUsage uint64
}

// cadvisorClient represent fixed client of cadvisor
type CadvisorClient struct {
	uri string
	staticClient *client.Client
}


// NewCadvisor return an fixed CadvisorClient instance
func NewCadvisor(uri string)(*CadvisorClient, error){
	staticClient,err := client.NewClient(buildCadvisorDia(uri))
	if err != nil{
		return nil, err
	}
	return &CadvisorClient{
		uri:          buildCadvisorDia(uri),
		staticClient: staticClient,
	}, nil
}

// MachineInfo get memory of machine
func (c *CadvisorClient)MachineInfo()(MachineInfo,error){
	mInfo,err := c.staticClient.MachineInfo()
	if err != nil{
		return MachineInfo{}, err
	}
	return MachineInfo{MemCapacity: mInfo.MemoryCapacity/1024}, nil
}

// ContainerInfo get memory information of container
func(c *CadvisorClient)ContainerInfo(containerName string)(ContainerInfo, error){
	// default return the latest record
	query := cadvisorv1.ContainerInfoRequest{NumStats: 1}

	cInfo,err := c.staticClient.DockerContainer(containerName, &query)
	if err != nil{
		return ContainerInfo{}, err
	}
	sCInfo := ContainerInfo{MemRss: 0, MemUsage: 0}
	for _,cs := range cInfo.Stats{
		// memory.Usage'value same as docker status show
		sCInfo.MemRss += cs.Memory.RSS / 1024
		sCInfo.MemUsage += cs.Memory.Usage / 1024
	}

	return sCInfo, nil
}

func buildCadvisorDia(url string)string{
	return "http://"+url+":4194"
}