package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	// CLI flags
	imagePtr := flag.String("image", "", "Docker image name to search for")
	kubeconfig := flag.String("kubeconfig", "/Users/foobar/.kube/config", "absolute path to the kubeconfig file")

	flag.Parse()

	if *imagePtr == "" {
		log.Fatal("Please provide an image name using --image flag")
	}

	if *kubeconfig == "" {
		log.Fatal("Please provide kubeconfig path using --kubeconfig flag. Ex: /home/user/.kube/config")
	}

	// Load config
	config, err := clientcmd.LoadFromFile(*kubeconfig)
	if err != nil {
		log.Fatalf("Failed to load kubeconfig: %v", err)
	}
	// Iterate through all contexts
	for contextName := range config.Contexts {
		fmt.Printf("Context: %s\n", contextName)

		// Load client for the context
		restConfig, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{CurrentContext: contextName},
		).ClientConfig()

		if err != nil {
			fmt.Printf("Failed to get client for context %s: %v\n", contextName, err)
			continue
		}

		// Create a Kubernetes client
		clientset, err := kubernetes.NewForConfig(restConfig)
		if err != nil {
			fmt.Printf("Failed to create clientset for context %s: %v\n", contextName, err)
			continue
		}

		// List all namespaces
		namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Printf("Failed to list namespaces for context %s: %v\n", contextName, err)
			continue
		}

		// Iterate through namespaces and print pods
		for _, ns := range namespaces.Items {
			pods, err := clientset.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				fmt.Printf("Failed to list pods in namespace %s for context %s: %v\n", ns.Name, contextName, err)
				continue
			}

			for _, pod := range pods.Items {
				// Check if the pod has the required label
				if chartLabel, exists := pod.Labels["chart"]; exists && strings.Contains(chartLabel, "service-") {
					// Check if any container has 'docker.io/busybox' image
					for _, container := range pod.Spec.Containers {
						// fmt.Printf("Pod: %s in Namespace: %s with image: %s\n", pod.Name, pod.Namespace, container.Image)
						if strings.Contains(container.Image, *imagePtr) {
							fmt.Printf("Context: %s | Namespace: %s | Pod: %s \n| Image: %s\n",
								contextName, ns.Name, pod.Name, container.Image)
							// fmt.Printf("Pod: %s in Namespace: %s matches the criteria\n", pod.Name, pod.Namespace)
							break // No need to check other containers in the same pod
						}
					}
				}
			}
		}
	}
}
