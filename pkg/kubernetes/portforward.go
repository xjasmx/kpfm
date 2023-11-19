package kubernetes

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"

	"github.com/xjasmx/kpfm/pkg/utils"
)

type PortForwarder struct {
	Clientset    *kubernetes.Clientset
	RestConfig   *rest.Config
	Namespace    string
	ResourceName string
	ResourceType string
	LocalPort    int
	RemotePort   int
}

func NewPortForwarder(clientset *kubernetes.Clientset, config *rest.Config, namespace, resource, portMapping string) (*PortForwarder, error) {
	localPort, remotePort, err := utils.ParsePortMapping(portMapping)
	if err != nil {
		return nil, fmt.Errorf("invalid port mapping format: %v", err)
	}

	resourceType, resourceName, err := parseResource(resource)
	if err != nil {
		return nil, err
	}

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

func establishPortForwarding(ctx context.Context, pf *PortForwarder, podName string) error {
	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward", pf.Namespace, podName)
	hostIP := strings.TrimLeft(pf.RestConfig.Host, "https://")

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

func parseResource(resource string) (string, string, error) {
	parts := strings.Split(resource, "/")
	if len(parts) != 2 {
		return parts[0], "pod", nil
	}
	return parts[0], parts[1], nil
}
