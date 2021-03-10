package pod

import (
	"context"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listercorev1 "k8s.io/client-go/listers/core/v1"
)

// podController represent an monitor instance for pod
type podController struct {
	kubeClient kubernetes.Interface
	podLister  listercorev1.PodLister
}

func (c *podController) OnAdd(obj interface{}) {
	panic("implement me")
}

func (c *podController) OnUpdate(oldObj, newObj interface{}) {
	panic("implement me")
}

func (c *podController) OnDelete(obj interface{}) {
	panic("implement me")
}

// NewPodController is an factory method that return an podController instance
func NewPodController(kubeClient kubernetes.Interface, kubeInformer informers.SharedInformerFactory) *podController {
	controller := &podController{
		kubeClient: kubeClient,
	}
	podInformer := kubeInformer.Core().V1().Pods()
	controller.podLister = podInformer.Lister()
	podInformer.Informer().AddEventHandler(controller)
	return controller
}

func (c *podController) Run(ctx context.Context) error {
	return nil
}
