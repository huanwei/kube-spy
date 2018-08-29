package main

import (
	"flag"
	"github.com/golang/glog"
	"github.com/huanwei/kube-spy/pkg/spy"
)

//spy program entrypoint

func main() {
	// Command line arguments
	var (
		kubeconfig string
	)
	flag.StringVar(&kubeconfig, "kubeconfig", "/etc/kubernetes/kubelet.conf", "absolute path to the kubeconfig file")
	flag.Parse()

	defer glog.Flush()

	// Configure k8s API client and get spy config
	clientset := spy.GetClientset(kubeconfig)
	spyConfig := spy.GetConfig()

	// Get services
	services := spy.GetServices(clientset, spyConfig)

	// Connect to DB
	spy.ConnectDB(clientset, spyConfig)

	// Close connection when exit
	defer spy.DBClient.Close()

	// Close all previous chaos
	spy.CloseChaos(clientset, spyConfig)

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
			// No chaos for this service, skip
			if len(spyConfig.VictimServices[i].ChaosList) == 0 {
				continue
			}
			// Detect network environment
			cidrs, podNames := spy.GetPodsInfo(spy.GetPods(clientset, services[i]))
			spy.PingPods(services[i].Name, services[i].Namespace, nil, podNames, cidrs)
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
				cidrs, podNames := spy.GetPodsInfo(spy.GetPods(clientset, services[i]))
				spy.PingPods(services[i].Name, services[i].Namespace, &chaos, podNames, cidrs)
				// Clear chaos
				spy.ClearChaos(clientset, spyConfig)
			}
			// Detect network environment
			cidrs, podNames = spy.GetPodsInfo(spy.GetPods(clientset, services[i]))
			spy.PingPods(services[i].Name, services[i].Namespace, nil, podNames, cidrs)
		}

	}

}
