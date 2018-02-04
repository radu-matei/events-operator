/*
github.com/radu-matei/event-operator
*/package fake

import (
	v1alpha1 "github.com/radu-matei/events-operator/pkg/apis/eventprovider/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeEventProviders implements EventProviderInterface
type FakeEventProviders struct {
	Fake *FakeEventproviderV1alpha1
	ns   string
}

var eventprovidersResource = schema.GroupVersionResource{Group: "eventprovider.radu-matei.com", Version: "v1alpha1", Resource: "eventproviders"}

var eventprovidersKind = schema.GroupVersionKind{Group: "eventprovider.radu-matei.com", Version: "v1alpha1", Kind: "EventProvider"}

// Get takes name of the eventProvider, and returns the corresponding eventProvider object, and an error if there is any.
func (c *FakeEventProviders) Get(name string, options v1.GetOptions) (result *v1alpha1.EventProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(eventprovidersResource, c.ns, name), &v1alpha1.EventProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventProvider), err
}

// List takes label and field selectors, and returns the list of EventProviders that match those selectors.
func (c *FakeEventProviders) List(opts v1.ListOptions) (result *v1alpha1.EventProviderList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(eventprovidersResource, eventprovidersKind, c.ns, opts), &v1alpha1.EventProviderList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.EventProviderList{}
	for _, item := range obj.(*v1alpha1.EventProviderList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested eventProviders.
func (c *FakeEventProviders) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(eventprovidersResource, c.ns, opts))

}

// Create takes the representation of a eventProvider and creates it.  Returns the server's representation of the eventProvider, and an error, if there is any.
func (c *FakeEventProviders) Create(eventProvider *v1alpha1.EventProvider) (result *v1alpha1.EventProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(eventprovidersResource, c.ns, eventProvider), &v1alpha1.EventProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventProvider), err
}

// Update takes the representation of a eventProvider and updates it. Returns the server's representation of the eventProvider, and an error, if there is any.
func (c *FakeEventProviders) Update(eventProvider *v1alpha1.EventProvider) (result *v1alpha1.EventProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(eventprovidersResource, c.ns, eventProvider), &v1alpha1.EventProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventProvider), err
}

// Delete takes name of the eventProvider and deletes it. Returns an error if one occurs.
func (c *FakeEventProviders) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(eventprovidersResource, c.ns, name), &v1alpha1.EventProvider{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeEventProviders) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(eventprovidersResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.EventProviderList{})
	return err
}

// Patch applies the patch and returns the patched eventProvider.
func (c *FakeEventProviders) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.EventProvider, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(eventprovidersResource, c.ns, name, data, subresources...), &v1alpha1.EventProvider{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.EventProvider), err
}
