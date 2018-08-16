package main

import (
	"fmt"
	"flag"
	"gopkg.in/resty.v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	"time"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"bytes"
	"github.com/huanwei/kube-spy/pkg/spy"
)
//spy program entrypoint


func main(){
	// Command line arguments
	var (
		kubeconfig string
	)
	flag.StringVar(&kubeconfig, "kubeconfig", "/etc/kubernetes/kubelet.conf", "absolute path to the kubeconfig file")
	flag.Parse()

	// Uses the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	// Creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// Read config file
	configFile,err:=ioutil.ReadFile("/spy/spy.conf")
	decoder:=yaml.NewDecoder(bytes.NewReader(configFile))

	// Decode config file
	spyConfig:=spy.Config{}
	decoder.Decode(&spyConfig)
	glog.Infof("decoded:\n%v",spyConfig)
	glog.Flush()

	// Get services
	services:=make([]*v1.Service,len(spyConfig.Service_list))
	for i,serviceName :=range spyConfig.Service_list{
		services[i],err=clientset.CoreV1().Services(spyConfig.Namespace).Get(serviceName,meta_v1.GetOptions{})
		if err!=nil{
			glog.Errorf("Fail to get service %s : %s",serviceName,err)
		}
	}
	glog.Flush()

	var host string
	glog.Infof("API service type: %s",services[0].Spec.Type)
	if services[0].Spec.Type==v1.ServiceType("ClusterIP"){
		host=services[0].Spec.ClusterIP
		glog.Infof("Service clusterIP: %s",host)

	}else {
		//TODO: other service type
	}

	glog.Flush()

	resp, err := resty.R().Get("http://"+host+"/get")
	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Received At: %v", resp.ReceivedAt())
	fmt.Printf("\nResponse Body: %v", resp)     // or resp.String() or string(resp.Body())
	glog.Infof("\nResponse Body: %v", resp)

	glog.Flush()
	// Wait for terminating
	for {
		time.Sleep(time.Duration(10) * time.Second)
	}


}