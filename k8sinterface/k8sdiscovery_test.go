package k8sinterface

import (
	"testing"

	"github.com/stretchr/testify/assert"

	wlidpkg "github.com/armosec/utils-k8s-go/wlid"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func TestResourceGroupToString(t *testing.T) {
	InitializeMapResourcesMock()

	allResources := ResourceGroupToString("*", "*", "*")
	if len(allResources) != len(ResourceGroupMapping) {
		t.Errorf("Expected len: %d, received: %d", len(ResourceGroupMapping), len(allResources))
	}
	pod := ResourceGroupToString("*", "*", "Pod")
	if len(pod) == 0 || pod[0] != "/v1/pods" {
		t.Errorf("pod: %v", pod)
	}
	deployments := ResourceGroupToString("*", "*", "Deployment")
	if len(deployments) == 0 || deployments[0] != "apps/v1/deployments" {
		t.Errorf("deployments: %v", deployments)
	}
	cronjobs := ResourceGroupToString("*", "*", "cronjobs")
	if len(cronjobs) == 0 || cronjobs[0] != "batch/v1beta1/cronjobs" {
		t.Errorf("cronjobs: %v", cronjobs)
	}
}

func TestGetGroupVersionResource(t *testing.T) {
	InitializeMapResourcesMock()
	wlid := "wlid://cluster-david-v1/namespace-default/deployment-nginx-deployment"
	r, err := GetGroupVersionResource(wlidpkg.GetKindFromWlid(wlid))
	if err != nil {
		t.Error(err)
		return
	}
	if r.Group != "apps" {
		t.Errorf("wrong group")
	}
	if r.Version != "v1" {
		t.Errorf("wrong Version")
	}
	if r.Resource != "deployments" {
		t.Errorf("wrong Resource")
	}

	r2, err := GetGroupVersionResource("NetworkPolicy")
	if err != nil {
		t.Error(err)
		return
	}
	if r2.Resource != "networkpolicies" {
		t.Errorf("wrong Resource")
	}
}

func TestIsNamespaceScope(t *testing.T) {
	InitializeMapResourcesMock()
	ResourceClusterScope = ResourceClusterScopeMock
	assert.True(t, IsNamespaceScope(&schema.GroupVersionResource{Group: "apps", Version: "v1", Resource: "deployments"}))
	assert.True(t, IsNamespaceScope(&schema.GroupVersionResource{Version: "v1", Resource: "serviceaccounts"}))
	assert.False(t, IsNamespaceScope(&schema.GroupVersionResource{Version: "v1", Resource: "nodes"}))
}