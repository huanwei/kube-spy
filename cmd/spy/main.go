package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/huanwei/kube-spy/pkg/spy"
	client_v2 "github.com/influxdata/influxdb/client/v2"
	"log"
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

	c, err := client_v2.NewHTTPClient(client_v2.HTTPConfig{
		Addr:     "http://192.168.102.238:8086",
		Username: "kubespy",
		Password: "kubespy",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Get services
	services := spy.GetServices(clientset, spyConfig)

	var host string
	// Get API server address
	if spyConfig.APIServerAddr == "" {
		host = spy.GetHost(clientset, services[0])
	} else {
		host = spyConfig.APIServerAddr
	}

	glog.Infof("There are %d chaos, %d test case in the list", len(spyConfig.ChaosList), len(spyConfig.TestCaseList))


	glog.Infof("len(services)=%d",len(services))
	// Len(chaos) + 1 tests, first one as normal test
	for i := -1; i < len(services); i++ {
		glog.Info("test ",i)
		if i == -1 {
			// Normal test
			glog.Infof("Normal test")
			spy.Dotests(spyConfig, host)
		} else {
			// Chaos tests
			for _, chaos := range spyConfig.ChaosList {
				glog.Infof("Chaos test: %v", chaos)
				// Add chaos
				err := spy.AddChaos(clientset, spyConfig, services[i], &chaos)
				if err != nil {
					glog.Errorf("Adding chaos error: %s", err)
				}
				// Start test
				spy.Dotests(spyConfig, host)
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
