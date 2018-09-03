package pkg

import (
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func GetService(clientset *kubernetes.Clientset, serviceName string) *v1.Service {
	service, err := clientset.CoreV1().Services("default").Get(serviceName, meta_v1.GetOptions{})
	if err != nil {
		glog.Errorf("Fail to get service %s : %s", serviceName, err)
		glog.Flush()
		panic(err.Error())
	}
	return service
}

// todo: Here service clusterIP is not the Pod IP. And it can be "None" if the service is statefulset headless service.
func GetHost(clientset *kubernetes.Clientset, service *v1.Service) string {
	var host string
	glog.Infof("API service type: %s", service.Spec.Type)
	if service.Spec.Type == v1.ServiceType("ClusterIP") {
		host = service.Spec.ClusterIP
		glog.Infof("Service clusterIP: %s", host)
	} else {
		//TODO: other service types
		// todo: Here v1.ServiceTypeNodePort can be exactly the Pod IP.
		glog.Infof("not implemented.")
	}
	return host
}
