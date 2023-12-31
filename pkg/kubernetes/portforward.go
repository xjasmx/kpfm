package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	"github.com/xjasmx/kpfm/pkg/utils"
)

type Clientset interface {
	CoreV1() corev1.CoreV1Interface
}

type PortForwarder struct {
	Clientset    Clientset
	RestConfig   *rest.Config
	Namespace    string
	ResourceName string
	ResourceType string
	LocalPort    int
	RemotePort   int
}

func NewPortForwarder(clientset Clientset, config *rest.Config, namespace, resource, portMapping string) (*PortForwarder, error) {
	err := validateInputs(namespace, resource, portMapping)
	if err != nil {
		return nil, err
	}

	localPort, remotePort, err := utils.ParsePortMapping(portMapping)
	if err != nil {
		return nil, err
	}

	resourceType, resourceName := parseResource(resource)

	return &PortForwarder{
		Clientset:    clientset,
		RestConfig:   config,
		Namespace:    namespace,
		ResourceName: resourceName,
		ResourceType: resourceType,
		LocalPort:    localPort,
		RemotePort:   remotePort,
	}, nil
}

func (pf *PortForwarder) Start(ctx context.Context) error {
	podName, err := resolvePodName(ctx, pf)
	if err != nil {
		return fmt.Errorf("failed to resolve pod name: %w", err)
	}

	return establishPortForwarding(ctx, pf, podName)
}

func validateInputs(namespace, resource, portMapping string) error {
	if namespace == "" {
		return fmt.Errorf("namespace cannot be empty")
	}

	if resource == "" {
		return fmt.Errorf("resource cannot be empty")
	}

	if portMapping == "" {
		return fmt.Errorf("portMapping cannot be empty")
	}

	return nil
}

func establishPortForwarding(ctx context.Context, pf *PortForwarder, podName string) error {
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", pf.Namespace, podName)
	hostIP := strings.TrimPrefix(pf.RestConfig.Host, "https://")

	transport, upgrader, err := spdy.RoundTripperFor(pf.RestConfig)
	if err != nil {
		return fmt.Errorf("failed to create round tripper: %w", err)
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", &url.URL{Scheme: "https", Path: path, Host: hostIP})
	ports := []string{fmt.Sprintf("%d:%d", pf.LocalPort, pf.RemotePort)}
	pfwd, err := portforward.New(dialer, ports, ctx.Done(), make(chan struct{}), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to create portforwarder: %w", err)
	}

	fmt.Printf("Port forwarding successful: %s:%d -> %d\n", pf.ResourceName, pf.RemotePort, pf.LocalPort)

	return pfwd.ForwardPorts()
}

func resolvePodName(ctx context.Context, pf *PortForwarder) (string, error) {
	switch pf.ResourceType {
	case "pod":
		return pf.ResourceName, nil
	case "service":
		coreClient := pf.Clientset.CoreV1()
		svc, err := coreClient.Services(pf.Namespace).Get(ctx, pf.ResourceName, metav1.GetOptions{})
		if err != nil {
			return "", fmt.Errorf("error getting service: %w", err)
		}
		set := labels.Set(svc.Spec.Selector)
		pods, err := coreClient.Pods(pf.Namespace).List(ctx, metav1.ListOptions{LabelSelector: set.AsSelector().String()})
		if err != nil {
			return "", fmt.Errorf("error finding pods for service: %w", err)
		}
		if len(pods.Items) == 0 {
			return "", fmt.Errorf("no pods found for service: %s", pf.ResourceName)
		}
		return pods.Items[0].Name, nil
	default:
		return "", fmt.Errorf("unsupported resource type: %s", pf.ResourceType)
	}
}

func parseResource(resource string) (string, string) {
	parts := strings.Split(resource, "/")
	if len(parts) != 2 {
		return parts[0], "pod"
	}
	return parts[0], parts[1]
}
