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

	// Connect to DB
	spy.ConnectDB(clientset, spyConfig)

	var host string
	// Get API server address
	if spyConfig.APIServerAddr == "" {
		host = services[0].Spec.ClusterIP
	} else {
		host = spyConfig.APIServerAddr
	}

	glog.Infof("There are %d services, %d test cases in the list", len(spyConfig.VictimServices), len(spyConfig.TestCases))

	// Len(chaos) + 1 tests, first one as normal test
	for i := -1; i < len(services); i++ {
		if i == -1 {
			// Normal test
			glog.Infof("None chaos test")
			spy.Dotests(spyConfig, host, nil, nil)
		} else {
			if len(spyConfig.VictimServices[i].ChaosList) == 0 {
				continue
			}
			// Detect network environment
			cidrs, podNames := spy.GetPods(clientset, services[i])
			delay, loss := spy.PingPods(cidrs)
			spy.StorePingResults(services[i].Name, services[i].Namespace, nil, podNames, delay, loss)
			// Chaos tests
			for _, chaos := range spyConfig.VictimServices[i].ChaosList {
				glog.Infof("Chaos test: Victim %s, Chaos %v", spyConfig.VictimServices[i].Name, chaos)
				// Add chaos
				err := spy.AddChaos(clientset, spyConfig, services[i], &chaos)
				if err != nil {
					glog.Errorf("Adding chaos error: %s", err)
				}
				// Do tests
				spy.Dotests(spyConfig, host, &spyConfig.VictimServices[i], &chaos)
				// Detect network environment
				cidrs, podNames := spy.GetPods(clientset, services[i])
				delay, loss := spy.PingPods(cidrs)
				spy.StorePingResults(services[i].Name, services[i].Namespace, &chaos, podNames, delay, loss)
				// Clear chaos
				spy.ClearChaos(clientset, spyConfig)
			}
			// Detect network environment
			cidrs, podNames = spy.GetPods(clientset, services[i])
			delay, loss = spy.PingPods(cidrs)
			spy.StorePingResults(services[i].Name, services[i].Namespace, nil, podNames, delay, loss)
		}

	}

	glog.Flush()

	// Close connection when exit
	spy.DBClient.Close()
	// Wait for terminating
	for {
		time.Sleep(time.Duration(10) * time.Second)
	}

}
