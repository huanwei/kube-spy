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
	"fmt"
	"bufio"
	"strings"
	"github.com/huanwei/kube-chaos/pkg/exec"
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
	services := make([]*v1.Service, len(config.VictimServices))
	var err error
	// Get services
	for i, service := range config.VictimServices {
		services[i], err = clientset.CoreV1().Services(config.Namespace).Get(service.Name, meta_v1.GetOptions{})
		if err != nil {
			glog.Errorf("Fail to get service %s : %s", service.Name, err)
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
		glog.Warningf("Unsupported service type: %v", service.Spec.Type)

	}
	return host
}

func GetPods(clientset *kubernetes.Clientset, service *v1.Service) (cidrs, podNames []string) {
	selector := service.Spec.Selector

	pods ,err := clientset.CoreV1().Pods("").List(meta_v1.ListOptions{LabelSelector:"app="+selector["app"]})
	if err != nil{
		glog.Errorf(fmt.Sprintf("Failed to get pods:%s",err))
		return cidrs,podNames
	}
	for _,pod := range pods.Items{
		cidrs = append(cidrs, pod.Status.PodIP)
		podNames = append(podNames, pod.Name)
	}
	return cidrs,podNames
}

func PingPods(cidrs []string)  (delay,loss[]string){
	for _,cidr := range cidrs{
		e := exec.New()
		glog.Infof(fmt.Sprintf("ping "+cidr))
		data,err := e.Command("ping","-i","0.01","-c","100",cidr).CombinedOutput()
		if err!= nil{
			glog.Errorf(fmt.Sprintf("Failed to ping %s:%s",cidr,err))
		} else {
			scanner := bufio.NewScanner(bytes.NewBuffer(data))
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if len(line) == 0 {
					continue
				}
				if strings.Contains(line, "transmitted") {
					glog.Infof(fmt.Sprintf("%s",line))
					parts := strings.Split(line, " ")
					percent := strings.Split(parts[5],"!")[0]
					loss = append(loss, percent)
				}
				if strings.Contains(line, "rtt") {
					glog.Infof(fmt.Sprintf("%s",line))
					delay = append(delay, line)
				}
			}
		}
	}
	return delay,loss
}

func StorePingResults(serviceName,namespace string, replicas int, ingress,egress string,podNames,delay,loss []string){
	for i := 0; i < len(delay); i++{
		AddPingResult(serviceName,namespace,replicas,ingress,egress,podNames[i],delay[i],loss[i])
	}
	SendPingResults()
}