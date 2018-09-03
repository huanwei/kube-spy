package main

import (
	"flag"
	"fmt"
	"http-test/pkg"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	targetPort := 8888
	port := 80
	var (
		kubeconfig  string
		nextService string
		service     string
	)
	flag.StringVar(&kubeconfig, "kubeconfig", "/etc/kubernetes/kubelet.conf", "absolute path to the kubeconfig file")
	flag.StringVar(&nextService, "nextService", "", "The http request destination service")
	flag.StringVar(&service, "service", "", "Current service")
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

	requestConfig := pkg.RequestConfig{
		Method: "GET",
		URL:    "/",
	}

	s := pkg.Server{}
	s.Initialize()
	if nextService != "" {
		host := pkg.GetHost(clientset, pkg.GetService(clientset, nextService))
		s.AddSendToNextHandler(requestConfig, host, service, nextService,port)
	} else {
		s.AddResponseHandler(requestConfig,service)
	}

	s.ListenAndServe(targetPort)
}
