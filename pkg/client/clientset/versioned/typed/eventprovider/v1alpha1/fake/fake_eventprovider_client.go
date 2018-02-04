/*
github.com/radu-matei/event-operator
*/package fake

import (
	v1alpha1 "github.com/radu-matei/events-operator/pkg/client/clientset/versioned/typed/eventprovider/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeEventproviderV1alpha1 struct {
	*testing.Fake
}

func (c *FakeEventproviderV1alpha1) EventProviders(namespace string) v1alpha1.EventProviderInterface {
	return &FakeEventProviders{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeEventproviderV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
