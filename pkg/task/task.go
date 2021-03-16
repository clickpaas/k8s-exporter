package task

import (
	"clickpaas-exporter/pkg/storage"
	"clickpaas-exporter/pkg/util"
	"context"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

const (
	DefaultAllocatedIfNodeSet = 1024*1024 * 2
)


// nodeInformationTask
type nodeInformationTask struct {
	kubeClient kubernetes.Interface
	cacheStorage *storage.CacheStorage
}

//podInformation information through k8s api-server
type podInformation struct {
	podName string
	containerId []string
	hasLimit bool
	MemLimit int64
	MemRss int64
}

// k8sNodeInformation
type k8sNodeInformation struct {
	nodeName string
	memCapacity int64			// kb
	memAllocatable int64		// kb
	ipAddress string

}


//NewNodeInformationTask return instance an nodeInformationTask
func NewNodeInformationTask(kubeClient kubernetes.Interface,  cacheStorage *storage.CacheStorage)*nodeInformationTask {
	task := &nodeInformationTask{
		kubeClient: kubeClient,
		cacheStorage: cacheStorage,
	}

	return task
}

// RunTask the entrypoint of task
func(t *nodeInformationTask)RunTask(){

	// get all node information
	allK8sNode,err := t.getAllNodeMachineryInformation()
	if err != nil{
		logrus.WithField("component", "task").Errorf("list k8s node failed, %s\n", err)
		return
	}
	allPodsInfo,err := t.getAllPodInformation()
	if err != nil{
		logrus.WithField("componet", "task").Errorf("list all pod throuh api-server failed %s\n", err)
		return
	}
	for _, node := range allK8sNode{
		cadvisorClient,err := util.NewCadvisor(node.ipAddress)
		if err != nil{
			logrus.WithField("component", "task").Errorf("get cadvisor client faild, %s\n",err)
			continue
		}
		machineInfo,err := cadvisorClient.MachineInfo()
		if err != nil{
			logrus.WithField("component", "task").Errorf("get machine info through cadvisor %s failed, %s\n", node.ipAddress, err)
			continue
		}
		memoryRecord := storage.MemoryRecord{
			// total memory
			Total: int64(machineInfo.MemCapacity),
			Allocatable: node.memAllocatable,
			Capacity: node.memCapacity,
			Allocated: 0,
		}
		allRelativePod := allPodsInfo[node.nodeName]
		for _,pod := range allRelativePod{
			// allocated
			memoryRecord.Allocated += t.getPodMemAllocatedOrUsage(pod, cadvisorClient, node.nodeName)
		}
		memoryRecord.Allocatable = memoryRecord.Total - memoryRecord.Allocated
		if memoryRecord.Allocatable < 0 {
			memoryRecord.Allocatable = 0
		}
		t.cacheStorage.Update(node.nodeName, memoryRecord)

	}
}

// getAllNodeMachineryInformation
func (t *nodeInformationTask)getAllNodeMachineryInformation()([]k8sNodeInformation,error){
	retK8Node := []k8sNodeInformation{}

	nodes,err := t.kubeClient.CoreV1().Nodes().List(context.TODO(),
		metav1.ListOptions{LabelSelector: labels.SelectorFromSet(map[string]string{"service": "USER_LABEL_VALUE", "TESTENV": "HOSTENVTYPE"}).String()})
	if err != nil{
		return nil, err
	}
	for _, node  := range nodes.Items{

		// if node is in Unschedulable ,then skip this node
		if node.Spec.Unschedulable {
			continue
		}
		var nodeInternalIp string
		for _,adr := range node.Status.Addresses{
			if adr.Type == "InternalIP"{
				nodeInternalIp = adr.Address
			}
		}
		retK8Node = append(retK8Node, k8sNodeInformation{
			nodeName:       node.Name,
			ipAddress: nodeInternalIp,
			memCapacity:    node.Status.Capacity.Memory().Value() / 1024,
			memAllocatable: node.Status.Allocatable.Memory().Value() / 1024,
		})
	}

	return retK8Node, nil
}

// getAllPodInformation get all pod information through k8s api-server
func (t *nodeInformationTask)getAllPodInformation()(map[string][]podInformation, error){
	retPodInfo := make(map[string][]podInformation)
	plist,err := t.kubeClient.CoreV1().Pods(v1.NamespaceAll).List(context.TODO(), metav1.ListOptions{})
	if err != nil{
		return nil, err
	}
	for _,pod := range plist.Items{
		podinfo := podInformation{hasLimit: false, podName: pod.Name}
		cNameIdMapping := getContainerNameIpMapping(&pod)
		// any given pod may have one more containers ,so here
		for _,container := range pod.Spec.Containers{
			if !container.Resources.Limits.Memory().IsZero() {
				podinfo.hasLimit = true
			}
			if podinfo.hasLimit {
				podinfo.MemLimit += container.Resources.Limits.Memory().Value() / 1024
			}
			podinfo.containerId = append(podinfo.containerId, cNameIdMapping[container.Name])
		}
		// add pod into special node,that the pod located
		if _,ok := retPodInfo[pod.Spec.NodeName]; !ok {
			retPodInfo[pod.Spec.NodeName] = []podInformation{podinfo}
		} else {
			retPodInfo[pod.Spec.NodeName] = append(retPodInfo[pod.Spec.NodeName], podinfo)
		}

	}
	return retPodInfo, err
}


//getPodMemAllocatedOrRss get pod resource limit ,if resource limit existed, or use the pod's rss replace
// if get pod's rss failed, then use default allocated value replace, the default allocated value is 2GB
func(t *nodeInformationTask)getPodMemAllocatedOrUsage(information podInformation, cadVisorClient *util.CadvisorClient, nodename string)int64{

	if information.hasLimit{
		return information.MemLimit
	}
	var memUsed int64 = 0
	for _, container := range information.containerId{

		containerInfo,err := cadVisorClient.ContainerInfo(container)

		if err != nil{
			memUsed += DefaultAllocatedIfNodeSet
		}else {
			memUsed += int64(containerInfo.MemUsage)
		}
	}
	return memUsed
}


