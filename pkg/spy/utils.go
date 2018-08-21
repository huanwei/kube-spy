package spy

import (
	"bytes"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func GetConfig() *Config {
	// Read config file
	configFile, err := ioutil.ReadFile("/spy/spy.conf")
	if err != nil {
		glog.Errorf("Fail to open spy config : %v", err)
		glog.Flush()
		panic(err.Error())
	}
	decoder := yaml.NewDecoder(bytes.NewReader(configFile))

	// Decode config file
	spyConfig := &Config{}
	decoder.Decode(spyConfig)
	glog.Infof("decoded:\n%v", spyConfig)
	glog.Flush()

	return spyConfig
}

func GetClientset(kubeconfig string) *kubernetes.Clientset {
	// Uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		glog.Errorf("Fail to build config from flags : %v", err)
		glog.Flush()
		panic(err.Error())
	}

	// Creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Errorf("Fail to create the clientset : %v", err)
		glog.Flush()
		panic(err.Error())
	}

	return clientset
}

func GetServices(clientset *kubernetes.Clientset, config *Config) []*v1.Service {
	// Create service array
	services := make([]*v1.Service, len(config.ServiceList))
	var err error
	// Get services
	for i, serviceName := range config.ServiceList {
		services[i], err = clientset.CoreV1().Services(config.Namespace).Get(serviceName, meta_v1.GetOptions{})
		if err != nil {
			glog.Errorf("Fail to get service %s : %s", serviceName, err)
			glog.Flush()
			panic(err.Error())
		}
	}
	return services
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
		glog.Warningf("Unsupported service type: %v",service.Spec.Type)

	}
	return host
}
