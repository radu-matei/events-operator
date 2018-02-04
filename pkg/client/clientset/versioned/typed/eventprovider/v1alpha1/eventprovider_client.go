/*
github.com/radu-matei/event-operator
*/package v1alpha1

import (
	v1alpha1 "github.com/radu-matei/events-operator/pkg/apis/eventprovider/v1alpha1"
	"github.com/radu-matei/events-operator/pkg/client/clientset/versioned/scheme"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type EventproviderV1alpha1Interface interface {
	RESTClient() rest.Interface
	EventProvidersGetter
}

// EventproviderV1alpha1Client is used to interact with features provided by the eventprovider.radu-matei.com group.
type EventproviderV1alpha1Client struct {
	restClient rest.Interface
}

func (c *EventproviderV1alpha1Client) EventProviders(namespace string) EventProviderInterface {
	return newEventProviders(c, namespace)
}

// NewForConfig creates a new EventproviderV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*EventproviderV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &EventproviderV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new EventproviderV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *EventproviderV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new EventproviderV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *EventproviderV1alpha1Client {
	return &EventproviderV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: scheme.Codecs}

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *EventproviderV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
