/*
github.com/radu-matei/event-operator
*/
// This file was automatically generated by informer-gen

package v1alpha1

import (
	internalinterfaces "github.com/radu-matei/events-operator/pkg/client/informers/externalversions/internalinterfaces"
)

// Interface provides access to all the informers in this group version.
type Interface interface {
	// EventProviders returns a EventProviderInformer.
	EventProviders() EventProviderInformer
}

type version struct {
	factory          internalinterfaces.SharedInformerFactory
	namespace        string
	tweakListOptions internalinterfaces.TweakListOptionsFunc
}

// New returns a new Interface.
func New(f internalinterfaces.SharedInformerFactory, namespace string, tweakListOptions internalinterfaces.TweakListOptionsFunc) Interface {
	return &version{factory: f, namespace: namespace, tweakListOptions: tweakListOptions}
}

// EventProviders returns a EventProviderInformer.
func (v *version) EventProviders() EventProviderInformer {
	return &eventProviderInformer{factory: v.factory, namespace: v.namespace, tweakListOptions: v.tweakListOptions}
}