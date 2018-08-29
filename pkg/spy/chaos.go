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
	pods := GetPods(clientset, service)

	// Control replicas
	ChangeReplicas(clientset, pods, int32(chaos.Replica), config.Namespace)

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

	// Able to select some of the pods to do chaos
	chaosPods := GetPartPods(clientset, service, chaos.Range)

	// Open these pods' chaos
	for _, pod := range chaosPods {
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
		pods := GetPods(clientset, service)

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
		time.Sleep(time.Duration(50 * time.Millisecond))
	}
	return nil
}

// Clear chaos in the previous influenced pods
func ClearChaos(clientset *kubernetes.Clientset, config *Config) {
	// Find previous chaos pods
	if PodsInChaos != nil {
		glog.Infof("Clear previous chaos: %s...%d pod(s)", PodsInChaos[0], len(PodsInChaos))
	} else {
		glog.Infof("No previous chaos to clear")
	}
	var (
		err error
		pod *v1.Pod
	)

	// If any
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
			time.Sleep(time.Duration(50 * time.Millisecond))
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
		glog.V(3).Infof("Clearing chaos on node \"%s\"...", node.Name)
		newAnnotations := node.Annotations
		newAnnotations["kubernetes.io/clear-chaos"] = " "
		node.SetAnnotations(newAnnotations)
		_, err := clientset.CoreV1().Nodes().UpdateStatus(node.DeepCopy())
		if err != nil {
			return errors.New(fmt.Sprintf("fail to update node status %s : %s", node.Name, err))
		}
	}

	// Wait for all node's chaos to close
	cnt := 1
	for {
		nodes, err = clientset.CoreV1().Nodes().List(meta_v1.ListOptions{LabelSelector: "chaos=on"})
		if err != nil {
			return errors.New(fmt.Sprintf("fail to list chaos nodes : %s", err))
		}

		if len(nodes.Items) == 0 {
			break
		}
		glog.V(3).Infof("Check nodes' chaos, try no. %d", cnt)
		cnt++
		time.Sleep(100 * time.Millisecond)
	}

	glog.V(3).Infof("Chaos cleared")

	return nil
}

// Control replicas via their deployment
func ChangeReplicas(clientset *kubernetes.Clientset, pods *v1.PodList, replica int32, namespace string) {
	if replica == 0 {
		return
	}
	for _, cref := range pods.Items[0].OwnerReferences {
		if !*cref.Controller {
			continue
		}
		replicaSet, err := clientset.AppsV1().ReplicaSets(namespace).Get(cref.Name, meta_v1.GetOptions{})
		if err != nil {
			glog.Errorf("Fail to find ReplicaSet %s: %s", cref.Name, err)
			continue
		}
		for _, dref := range replicaSet.OwnerReferences {
			if !*dref.Controller {
				continue
			}
			deployment, err := clientset.AppsV1().Deployments(namespace).Get(dref.Name, meta_v1.GetOptions{})
			if err != nil {
				glog.Errorf("Fail to find Deploymnet %s: %s", cref.Name, err)
				continue
			}
			glog.Infof("Previous replicas: %d", *deployment.Spec.Replicas)
			deployment.Spec.Replicas = &replica

			_, err = clientset.AppsV1().Deployments(namespace).Update(deployment.DeepCopy())
			if err != nil {
				glog.Errorf("Scale error: %s", err)
				continue
			}
			deployment, err = clientset.AppsV1().Deployments(namespace).Get(dref.Name, meta_v1.GetOptions{})
			if err != nil {
				glog.Errorf("Fail to find Deploymnet %s: %s", cref.Name, err)
				continue
			}
			glog.Infof("Deploymnet %s scaled to %d", deployment.Name, *deployment.Spec.Replicas)
		}
	}

}
