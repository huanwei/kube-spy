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
)
//spy program entrypoint


func main(){
	// Command line arguments
	var kubeconfig string
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

	// For testing
	namespace:="default"
	serviceNames:=[]string{"hello-1","hello-2"}

	services:=make([]*v1.Service,len(serviceNames))
	for i,serviceName :=range serviceNames{
		services[i],err=clientset.CoreV1().Services(namespace).Get(serviceNames[i],meta_v1.GetOptions{})
		if err!=nil{
			glog.Errorf("Fail to get service %s : %s",serviceName,err)
		}
	}
	glog.Flush()

	resp, err := resty.R().Get("http://httpbin.org/get")


	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", resp.StatusCode())
	fmt.Printf("\nResponse Status: %v", resp.Status())
	fmt.Printf("\nResponse Time: %v", resp.Time())
	fmt.Printf("\nResponse Received At: %v", resp.ReceivedAt())
	fmt.Printf("\nResponse Body: %v", resp)     // or resp.String() or string(resp.Body())

	// Wait for terminating
	for {
		time.Sleep(time.Duration(10) * time.Second)
	}


}