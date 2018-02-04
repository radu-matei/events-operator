/*
github.com/radu-matei/event-operator
*/package v1alpha1

import (
	v1alpha1 "github.com/radu-matei/events-operator/pkg/apis/eventprovider/v1alpha1"
	scheme "github.com/radu-matei/events-operator/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// EventProvidersGetter has a method to return a EventProviderInterface.
// A group's client should implement this interface.
type EventProvidersGetter interface {
	EventProviders(namespace string) EventProviderInterface
}

// EventProviderInterface has methods to work with EventProvider resources.
type EventProviderInterface interface {
	Create(*v1alpha1.EventProvider) (*v1alpha1.EventProvider, error)
	Update(*v1alpha1.EventProvider) (*v1alpha1.EventProvider, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.EventProvider, error)
	List(opts v1.ListOptions) (*v1alpha1.EventProviderList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.EventProvider, err error)
	EventProviderExpansion
}

// eventProviders implements EventProviderInterface
type eventProviders struct {
	client rest.Interface
	ns     string
}

// newEventProviders returns a EventProviders
func newEventProviders(c *EventproviderV1alpha1Client, namespace string) *eventProviders {
	return &eventProviders{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the eventProvider, and returns the corresponding eventProvider object, and an error if there is any.
func (c *eventProviders) Get(name string, options v1.GetOptions) (result *v1alpha1.EventProvider, err error) {
	result = &v1alpha1.EventProvider{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("eventproviders").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of EventProviders that match those selectors.
func (c *eventProviders) List(opts v1.ListOptions) (result *v1alpha1.EventProviderList, err error) {
	result = &v1alpha1.EventProviderList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("eventproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested eventProviders.
func (c *eventProviders) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("eventproviders").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a eventProvider and creates it.  Returns the server's representation of the eventProvider, and an error, if there is any.
func (c *eventProviders) Create(eventProvider *v1alpha1.EventProvider) (result *v1alpha1.EventProvider, err error) {
	result = &v1alpha1.EventProvider{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("eventproviders").
		Body(eventProvider).
		Do().
		Into(result)
	return
}

// Update takes the representation of a eventProvider and updates it. Returns the server's representation of the eventProvider, and an error, if there is any.
func (c *eventProviders) Update(eventProvider *v1alpha1.EventProvider) (result *v1alpha1.EventProvider, err error) {
	result = &v1alpha1.EventProvider{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("eventproviders").
		Name(eventProvider.Name).
		Body(eventProvider).
		Do().
		Into(result)
	return
}

// Delete takes name of the eventProvider and deletes it. Returns an error if one occurs.
func (c *eventProviders) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("eventproviders").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *eventProviders) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("eventproviders").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched eventProvider.
func (c *eventProviders) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.EventProvider, err error) {
	result = &v1alpha1.EventProvider{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("eventproviders").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
