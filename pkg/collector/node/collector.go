package node

import (
	"clickpaas-exporter/conf"
	"clickpaas-exporter/pkg/storage"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	NameSpace  = "k8sexporter"
	SubSystem  = "node"
	K8sCluster = "dev"
)

type nodeCollector struct {
	numDesc   int
	nodeCache *storage.CacheStorage
	// metric from cadvisor

	// nodeMachineMemory the total memory of physical node
	// this metric get from cadvisor
	nodeMemTotalMachine *prometheus.Desc

	// metric from k8s apiServer
	// nodeAllocatedMemory represent the allocated by k8s
	// if pod is not set resource'limit ,then use rss replace
	nodeMemAllocated *prometheus.Desc
	// capacity - allocated = allocatable
	nodeMemAllocatable *prometheus.Desc
	// nodeCapacity
	nodeMemCapacity *prometheus.Desc
	// allocated percent
	nodeMemAllocatedPercent *prometheus.Desc
	// allocatable percent
	nodeMemAllocatablePercent *prometheus.Desc

}

func NewNodeCollector(nodeCache *storage.CacheStorage, config *conf.Config) *nodeCollector {
	collector := &nodeCollector{nodeCache: nodeCache}
	collector.nodeMemTotalMachine = prometheus.NewDesc(
		prometheus.BuildFQName(NameSpace, SubSystem, "memory_total"), "machinery total memory", []string{"node"}, map[string]string{"k8scluster": config.Cluster})
	collector.nodeMemAllocated = prometheus.NewDesc(
		prometheus.BuildFQName(NameSpace, SubSystem, "memory_allocated"), "allocated by k8s", []string{"node"}, map[string]string{"k8scluster": config.Cluster})
	collector.nodeMemAllocatable = prometheus.NewDesc(
		prometheus.BuildFQName(NameSpace, SubSystem, "memory_allocatable"), "allocatable ", []string{"node"}, map[string]string{"k8scluster": config.Cluster})
	collector.nodeMemCapacity = prometheus.NewDesc(
		prometheus.BuildFQName(NameSpace, SubSystem, "memory_capacity"), "capacity", []string{"node"}, map[string]string{"k8scluster": config.Cluster})
	collector.nodeMemAllocatablePercent = prometheus.NewDesc(
		prometheus.BuildFQName(NameSpace, SubSystem, "memory_allocatable_percent"), "allOcateable/total", []string{"node"}, map[string]string{"k8scluster": config.Cluster})
	collector.nodeMemAllocatedPercent = prometheus.NewDesc(
		prometheus.BuildFQName(NameSpace, SubSystem, "memory_allocated_percent"), "allocated / total", []string{"node"}, map[string]string{"k8scluster": config.Cluster})
	return collector
}

func (n *nodeCollector) Describe(descs chan<- *prometheus.Desc) {
	for i := 0; i < n.numDesc; i++ {
		descs <- n.nodeMemTotalMachine
		descs <- n.nodeMemAllocated
		descs <- n.nodeMemCapacity
		descs <- n.nodeMemAllocatable
		descs <- n.nodeMemAllocatedPercent
		descs <- n.nodeMemAllocatablePercent
	}
}

func (n *nodeCollector) Collect(metrics chan<- prometheus.Metric) {
	nodesInfo := n.nodeCache.DeepCopy()
	var (
		allocatedPercent float64
		allocatablePercent float64
	)
	for nodeName, record := range nodesInfo {
		if record.Total == 0{
			allocatablePercent = 0
			allocatedPercent = 0
		} else {
			allocatablePercent = float64(record.Allocatable) / float64(record.Total)
			allocatedPercent = float64(record.Allocated) / float64(record.Total)
		}
		metrics <- prometheus.MustNewConstMetric(n.nodeMemTotalMachine, prometheus.GaugeValue, float64(record.Total), nodeName)
		metrics <- prometheus.MustNewConstMetric(n.nodeMemAllocated, prometheus.GaugeValue, float64(record.Allocated), nodeName)
		metrics <- prometheus.MustNewConstMetric(n.nodeMemAllocatable, prometheus.GaugeValue, float64(record.Allocatable), nodeName)
		metrics <- prometheus.MustNewConstMetric(n.nodeMemCapacity, prometheus.GaugeValue, float64(record.Capacity), nodeName)
		metrics <- prometheus.MustNewConstMetric(n.nodeMemAllocatablePercent, prometheus.GaugeValue, allocatablePercent, nodeName)
		metrics <- prometheus.MustNewConstMetric(n.nodeMemAllocatedPercent, prometheus.GaugeValue, allocatedPercent, nodeName)
		n.numDesc += 1
	}
}
