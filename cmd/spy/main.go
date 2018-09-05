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
	spy.CloseChaos(clientset)

	// Account testcase number
	testCaseNum := 0
	for _, testcaselist := range spyConfig.TestCaseLists {
		testCaseNum += len(testcaselist.TestCases)
	}
	glog.Infof("There are %d services, %d testcase lists, %d testcases in config", len(spyConfig.VictimServices), len(spyConfig.TestCaseLists), testCaseNum)

	// Len(chaos) + 1 tests, first one as normal test
	for i := -1; i < len(services); i++ {
		if i == -1 {
			// Normal test
			glog.Infof("None chaos test")
			spy.Dotests(clientset, spyConfig, nil, nil)
		} else {
			// No chaos for this service, skip
			if len(spyConfig.VictimServices[i].ChaosList) == 0 {
				continue
			}
			if spyConfig.VictimServices[i].PingTimeout <= 1 {
				spyConfig.VictimServices[i].PingTimeout = 1
			} else {
				spyConfig.VictimServices[i].PingTimeout = spyConfig.VictimServices[i].PingTimeout - 1
			}
			for _, chaos := range spyConfig.VictimServices[i].ChaosList {
				stop := make(chan bool, 1)
				complete := make(chan bool, 1)
				glog.Infof("Chaos test: Victim %s, Chaos %v", spyConfig.VictimServices[i].Name, chaos)

				// Control replicas
				previousReplica := spy.ChangeReplicas(clientset, services[i], int32(chaos.Replica), spyConfig.Namespace)

				// Get pods after changing replicas
				pods := spy.GetPods(clientset, services[i])

				// Detect network environment before adding chaos
				cidrs, podNames := spy.GetPodsInfo(pods)
				go spy.PingPods(services[i].Name, services[i].Namespace, podNames, cidrs, &chaos, stop, complete, spyConfig.VictimServices[i].PingTimeout)

				time.Sleep(time.Duration(1) * time.Second)

				// Add chaos
				err := spy.AddChaos(clientset, spyConfig, services[i], &chaos, pods)
				if err != nil {
					glog.Errorf("Adding chaos error: %s", err)
				}

				time.Sleep(time.Duration(spyConfig.VictimServices[i].PingTimeout) * time.Second)

				// Do API tests
				spy.Dotests(clientset, spyConfig, &spyConfig.VictimServices[i], &chaos)

				time.Sleep(time.Duration(spyConfig.VictimServices[i].PingTimeout) * time.Second)

				// Clear chaos
				spy.ClearChaos(clientset, spyConfig)

				time.Sleep(time.Duration(3) * time.Second)

				stop <- true
				<-complete

				// Restore replicas
				spy.ChangeReplicas(clientset, services[i], int32(previousReplica), spyConfig.Namespace)
			}
		}
	}
}
