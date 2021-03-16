package main

import (
	"clickpaas-exporter/pkg/collector"
	"clickpaas-exporter/pkg/collector/node"
	"clickpaas-exporter/pkg/storage"
	"clickpaas-exporter/pkg/task"
	"fmt"
	cronv3 "github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"net/http"
	_ "net/http/pprof"
)

var (
	masterUrl      = pflag.String("masterUrl", "", "address of k8s apiServer")
	kubeConfig     = pflag.String("kubeConfig", "", "path of kubeConfig")

	configPath = pflag.String("config", "", "the config file")
	k8scluster = pflag.String("k8scluster", "dev", "k8s cluster ")
	// special options
)

var (
	kubeClient kubernetes.Interface

	restConfig *rest.Config
)

func init() {

}

func main() {

	pflag.Parse()

	//crConfig := conf.NewConfig(*configPath)

	if err := buildKubeConfig(*masterUrl, *kubeConfig); err != nil {
		logrus.Panicf("create k8s config failed, %s", err)
	}
	var err error
	kubeClient, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		logrus.Panic("build ku8s client failed, %s", err)
	}

	cacheStorage := storage.NewCacheStorage()
	cacheStorage.Update("haha", storage.MemoryRecord{
		Total:       0,
		Allocated:   0,
		Capacity:    0,
		Allocatable: 0,
	})

	nodeInformationTask := task.NewNodeInformationTask(kubeClient, cacheStorage)

	// run crontab task
	// default every 15second gather all information, then restore template storage, waiting for prometheus collector gathering
	crontask := cronv3.New()
	crontask.AddFunc("@every 15s", nodeInformationTask.RunTask)
	crontask.Start()
	logrus.Infof("crontab start")

	collectorHandler := collector.NewCollectorHandler(false)
	collectorHandler.MustRegister(node.NewNodeCollector(cacheStorage, *k8scluster))

	mux := http.DefaultServeMux
	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "ok")
	})
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "hello world")
	})
	mux.Handle("/metrics", collectorHandler)
	logrus.Panic(http.ListenAndServe(":8989", mux))
}

// buildKubeConfig build kubeConfig from kube config file
// default kube config file is located at ~/.kube/config
func buildKubeConfig(masterUrl, kubeConfig string) (err error) {

	if kubeConfig != "" {
		restConfig, err = clientcmd.BuildConfigFromFlags(masterUrl, kubeConfig)
	} else {
		restConfig, err = rest.InClusterConfig()
	}
	return err
}
