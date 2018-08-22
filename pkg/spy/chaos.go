package spy

import (
	"errors"
	"fmt"
	"github.com/golang/glog"
	"k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

var PodsInChaos []string

// Add chaos to the service's pods
func AddChaos(clientset *kubernetes.Clientset, config *Config, service *v1.Service, chaos *Chaos) error {
	labelselector:=""
	for selectors,values:=range service.Spec.Selector{
		if labelselector==""{
			labelselector+=selectors+"="+values
		}else {
			labelselector+=","+selectors+"="+values
		}
	}
	glog.Infof("labelselector: %s",labelselector)

	//deploy,err:=clientset.AppsV1().Deployments(config.Namespace).Get(service.Name,meta_v1.GetOptions{})
	//if err!=nil{
	//	glog.Warningf("Deployment %s: not found: %s",service.Name,err)
	//}else {
	//	deploy.Status.Replicas=3
	//	clientset.AppsV1().Deployments(config.Namespace).UpdateStatus(deploy.DeepCopy())
	//}

	pods, err := clientset.CoreV1().Pods(config.Namespace).List(meta_v1.ListOptions{LabelSelector: labelselector})
	if err != nil {
		return errors.New(fmt.Sprintf("fail to list service %s's corresponding pods : %s", service.Name, err))
	}

	glog.Infof("%v",pods.Items[0].OwnerReferences)
	for _,ref:=range pods.Items[0].OwnerReferences{
		if *ref.Controller{
			replicaset,err:=clientset.AppsV1().ReplicaSets("default").Get(ref.Name,meta_v1.GetOptions{})
			if err!=nil{
				glog.Warningf("Fail to find ReplicaSet %s: %s",ref.Name,err)
			}else {
				replicaset.Status.Replicas=3
				clientset.AppsV1().ReplicaSets(config.Namespace).UpdateStatus(replicaset.DeepCopy())
			}
		}
	}

	// Find all nodes running this service's pods
	nodes := make(map[string]*v1.Node)
	for _, pod := range pods.Items {
		_, alreadyFound := nodes[pod.Spec.NodeName]
		if !alreadyFound {
			node, err := clientset.CoreV1().Nodes().Get(pod.Spec.NodeName, meta_v1.GetOptions{})
			if err != nil {
				return errors.New(fmt.Sprintf("fail to get node %s : %s", pod.Spec.NodeName, err))
			}
			nodes[pod.Spec.NodeName] = node
		}
	}
	glog.Infof("chaos nodes: ")
	// Open these nodes' chaos
	for name, node := range nodes {
		glog.Info(name)
		newLabels := node.Labels
		_, on := newLabels["chaos"]
		if on {
			continue
		}
		newLabels["chaos"] = "on"
		node.SetLabels(newLabels)
		_, err := clientset.CoreV1().Nodes().UpdateStatus(node.DeepCopy())
		if err != nil {
			return errors.New(fmt.Sprintf("fail to update node status %s : %s", node.Name, err))
		}
	}

	glog.Infof("Got pods:")
	// Open these pods' chaos
	for _, pod := range pods.Items {
		glog.Info(pod.Name)
		// Set labels
		newLabels := pod.Labels
		newLabels["chaos"] = "on"
		pod.SetLabels(newLabels)
		// Set annotations
		newAnnotations := pod.Annotations
		if newAnnotations == nil {
			newAnnotations = make(map[string]string)
		}
		newAnnotations["kubernetes.io/egress-chaos"] = chaos.Egress
		newAnnotations["kubernetes.io/done-egress-chaos"] = "no"
		newAnnotations["kubernetes.io/ingress-chaos"] = chaos.Ingress
		newAnnotations["kubernetes.io/done-ingress-chaos"] = "no"
		pod.SetAnnotations(newAnnotations)
		// Update pod

		_, err := clientset.CoreV1().Pods(config.Namespace).UpdateStatus(pod.DeepCopy())
		if err != nil {
			return errors.New(fmt.Sprintf("fail to update pod status %s : %s", pod.Name, err))
		}
		PodsInChaos = append(PodsInChaos, pod.Name)
	}

	// Wait for response
	for {
		allReady := true
		pods, err := clientset.CoreV1().Pods(config.Namespace).List(meta_v1.ListOptions{LabelSelector:labelselector})
		if err != nil {
			return errors.New(fmt.Sprintf("fail to list service %s's corresponding pods : %s", service.Name, err))
		}

		for _, pod := range pods.Items {
			done, _ := pod.Annotations["kubernetes.io/done-egress-chaos"]
			if done == "no" {
				allReady = false
				break
			}
		}
		if allReady {
			break
		}
		time.Sleep(time.Duration(50 * time.Microsecond))
	}
	return nil
}

// Clear chaos in the previous influenced pods
func ClearChaos(clientset *kubernetes.Clientset, config *Config) {
	if PodsInChaos != nil {
		glog.Infof("Clear previous chaos: %s...%d pod(s)", PodsInChaos[0], len(PodsInChaos))
	} else {
		glog.Infof("No previous chaos to clear")
	}
	var (
		err error
		pod *v1.Pod
	)
	if PodsInChaos != nil {
		for _, podName := range PodsInChaos {
			pod, err = clientset.CoreV1().Pods(config.Namespace).Get(podName, meta_v1.GetOptions{})
			if err != nil {
				glog.Errorf("Fail to get pod %s : %s", podName, err)
				continue
			}

			newAnnotations := pod.Annotations
			if newAnnotations == nil {
				newAnnotations = make(map[string]string)
			}
			newAnnotations["kubernetes.io/clear-ingress-chaos"] = "yes"
			newAnnotations["kubernetes.io/clear-egress-chaos"] = "yes"
			newAnnotations["kubernetes.io/done-ingress-chaos"] = "no"
			newAnnotations["kubernetes.io/done-egress-chaos"] = "no"
			pod.SetAnnotations(newAnnotations)
			_, err = clientset.CoreV1().Pods(config.Namespace).UpdateStatus(pod.DeepCopy())
			if err != nil {
				glog.Errorf("Fail to update pod status %s : %s", pod.Name, err)
				continue
			}
		}

		// Wait for response
		for {
			allReady := true
			for _, podName := range PodsInChaos {
				pod, err = clientset.CoreV1().Pods(config.Namespace).Get(podName, meta_v1.GetOptions{})
				if pod.Annotations == nil {
					continue
				}
				_, egressNotClear := pod.Annotations["kubernetes.io/egress-chaos"]
				_, ingressNotClear := pod.Annotations["kubernetes.io/ingress-chaos"]

				if egressNotClear || ingressNotClear {
					allReady = false
					break
				}
			}
			if allReady {
				break
			}
			time.Sleep(time.Duration(50 * time.Microsecond))
		}

		PodsInChaos = nil
	}
}

// Close all chaos nodes' chaos
func CloseChaos(clientset *kubernetes.Clientset, config *Config) error {
	// List all chaos nodes
	nodes, err := clientset.CoreV1().Nodes().List(meta_v1.ListOptions{LabelSelector: "chaos=on"})
	if err != nil {
		return errors.New(fmt.Sprintf("fail to list chaos nodes : %s", err))
	}

	// Set clear flag on nodes' annotation
	for _, node := range nodes.Items {
		newAnnotations := node.Annotations
		newAnnotations["kubernetes.io/clear-chaos"] = " "
		node.SetAnnotations(newAnnotations)
		_, err := clientset.CoreV1().Nodes().UpdateStatus(node.DeepCopy())
		if err != nil {
			return errors.New(fmt.Sprintf("fail to update node status %s : %s", node.Name, err))
		}
	}
	return nil
}
