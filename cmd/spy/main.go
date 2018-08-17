package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/huanwei/kube-spy/pkg/spy"
	"time"
)

//spy program entrypoint

func main() {
	// Command line arguments
	var (
		kubeconfig string
	)
	flag.StringVar(&kubeconfig, "kubeconfig", "/etc/kubernetes/kubelet.conf", "absolute path to the kubeconfig file")
	flag.Parse()

	// Configure k8s API client and get spy config
	clientset := spy.GetClientset(kubeconfig)
	spyConfig := spy.GetConfig()

	// Get services
	services := spy.GetServices(clientset, spyConfig)

	// Get API server address
	host := spy.GetHost(clientset, services[0])

	glog.Infof("There are %d chaos, %d test case in the list", len(spyConfig.ChaosList), len(spyConfig.TestCaseList))

	// Len(chaos) + 1 tests, first one as normal test
	for i := -1; i < len(services); i++ {
		if i == -1 {
			// Normal test
			glog.Infof("Normal test")
			spy.Dotests(spyConfig.TestCaseList, host)
		} else {
			// Chaos tests
			for _, chaos := range spyConfig.ChaosList {
				glog.Infof("Chaos test: %s", chaos)
				// Add chaos
				err := spy.AddChaos(clientset, spyConfig, services[i], chaos)
				if err != nil {
					glog.Errorf("Adding chaos error: %s", err)
				}
				// Start test
				spy.Dotests(spyConfig.TestCaseList, host)
				// Clear chaos
				spy.ClearChaos(clientset, spyConfig)
			}
		}

	}

	glog.Flush()

	// Wait for terminating
	for {
		time.Sleep(time.Duration(10) * time.Second)
	}

}
