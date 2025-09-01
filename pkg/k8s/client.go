package k8s

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// Client wraps the Kubernetes client
type Client struct {
	clientset *kubernetes.Clientset
	namespace string
}

// NewClient creates a new Kubernetes client
func NewClient(inCluster bool, configPath, namespace string) (*Client, error) {
	var config *rest.Config
	var err error

	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", configPath)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return &Client{
		clientset: clientset,
		namespace: namespace,
	}, nil
}

// GetVaultPods returns a list of Vault pods in the configured namespace
func (c *Client) GetVaultPods(ctx context.Context) ([]string, error) {
	pods, err := c.clientset.CoreV1().Pods(c.namespace).List(ctx, metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/instance=vault", // TODO Adjust label selector as needed (this one reads also the agent injector pod)
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list Vault pods: %w", err)
	}

	var podNames []string
	for _, pod := range pods.Items {
		if pod.Status.Phase == "Running" {
			podNames = append(podNames, pod.Name)
		}
	}

	return podNames, nil
}

// IsVaultReady checks if Vault pods are ready
func (c *Client) IsVaultReady(ctx context.Context) (bool, error) {
	pods, err := c.GetVaultPods(ctx)
	fmt.Println("Vault pods found:", pods)
	if err != nil {
		return false, err
	}

	return len(pods) > 0, nil
}
