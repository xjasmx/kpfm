package kubernetes

import (
	"testing"

	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
)

func TestNewPortForwarder(t *testing.T) {
	fakeClientset := fake.NewSimpleClientset()
	fakeConfig := &rest.Config{}

	tests := []struct {
		name        string
		namespace   string
		resource    string
		portMapping string
		expectErr   bool
	}{
		{
			name:        "Valid input",
			namespace:   "default",
			resource:    "service/my-service",
			portMapping: "8080:80",
			expectErr:   false,
		},
		{
			name:        "Valid input without resource type",
			namespace:   "default",
			resource:    "pod",
			portMapping: "8080:80",
			expectErr:   false,
		},
		{
			name:        "Invalid port mapping",
			namespace:   "default",
			resource:    "service/my-service",
			portMapping: "abc:def",
			expectErr:   true,
		},
		{
			name:        "Empty namespace",
			namespace:   "",
			resource:    "service/my-service",
			portMapping: "8080:80",
			expectErr:   true,
		},
		{
			name:        "Empty resource",
			namespace:   "default",
			resource:    "",
			portMapping: "8080:80",
			expectErr:   true,
		},
		{
			name:        "Empty port mapping",
			namespace:   "default",
			resource:    "service/my-service",
			portMapping: "",
			expectErr:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := NewPortForwarder(fakeClientset, fakeConfig, tc.namespace, tc.resource, tc.portMapping)
			if (err != nil) != tc.expectErr {
				t.Errorf("NewPortForwarder() error = %v, expectErr %v", err, tc.expectErr)
			}
		})
	}
}
