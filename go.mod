module clickpaas-exporter

go 1.16

require (
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b // indirect
	github.com/google/cadvisor v0.38.8
	github.com/prometheus/client_golang v1.9.0
	github.com/robfig/cron v1.2.0 // indirect
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.8.0
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	k8s.io/client-go v0.20.4
)
