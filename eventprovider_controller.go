package main

import (
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/golang/glog"
	"github.com/radu-matei/events-operator/pkg/apis/eventprovider/v1alpha1"
	clientset "github.com/radu-matei/events-operator/pkg/client/clientset/versioned"
	sscheme "github.com/radu-matei/events-operator/pkg/client/clientset/versioned/scheme"
	informers "github.com/radu-matei/events-operator/pkg/client/informers/externalversions"
	listers "github.com/radu-matei/events-operator/pkg/client/listers/eventprovider/v1alpha1"
	eventgrid "github.com/radu-matei/events-operator/pkg/eventgrid"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	kubeinformers "k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	typedcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	appslisters "k8s.io/client-go/listers/apps/v1"
	corelisters "k8s.io/client-go/listers/core/v1"
	extensionlisters "k8s.io/client-go/listers/extensions/v1beta1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
)

const controllerAgentName = "eventprovider_controller"

// Controller is the controller implementation for Foo resources
type Controller struct {
	kubeclientset kubernetes.Interface
	epclientset   clientset.Interface

	epLister listers.EventProviderLister
	epSynced cache.InformerSynced

	deploymentsLister appslisters.DeploymentLister
	deploymentsSynced cache.InformerSynced

	servicesLister corelisters.ServiceLister
	servicesSynced cache.InformerSynced

	ingressLister extensionlisters.IngressLister
	ingressSynced cache.InformerSynced

	queue    workqueue.RateLimitingInterface
	recorder record.EventRecorder
}

// NewController returns a new instance of a controller
func NewController(
	kubeclientset kubernetes.Interface,
	epclientset clientset.Interface,

	kubeInformerFactory kubeinformers.SharedInformerFactory,
	epInformerFactory informers.SharedInformerFactory) *Controller {

	epInformer := epInformerFactory.Eventprovider().V1alpha1().EventProviders()
	sscheme.AddToScheme(scheme.Scheme)

	deploymentInformer := kubeInformerFactory.Apps().V1().Deployments()
	serviceInformer := kubeInformerFactory.Core().V1().Services()
	ingressInformer := kubeInformerFactory.Extensions().V1beta1().Ingresses()

	glog.V(4).Info("Creating event broadcaster")
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&typedcorev1.EventSinkImpl{Interface: kubeclientset.CoreV1().Events("")})
	recorder := eventBroadcaster.NewRecorder(scheme.Scheme, corev1.EventSource{Component: controllerAgentName})

	c := &Controller{
		kubeclientset: kubeclientset,
		epclientset:   epclientset,

		epLister: epInformer.Lister(),
		epSynced: epInformer.Informer().HasSynced,

		deploymentsLister: deploymentInformer.Lister(),
		deploymentsSynced: deploymentInformer.Informer().HasSynced,

		servicesLister: serviceInformer.Lister(),
		servicesSynced: serviceInformer.Informer().HasSynced,

		ingressLister: ingressInformer.Lister(),
		ingressSynced: ingressInformer.Informer().HasSynced,

		queue:    workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "EventProviders"),
		recorder: recorder,
	}

	glog.Info("Setting up event handlers")
	// Set up an event handler for when EventProvider resources change
	epInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			glog.Info("AddFunc called with object: %v", obj)
			key, err := cache.MetaNamespaceKeyFunc(obj)
			if err == nil {
				c.queue.Add(key)
			}
		},
		UpdateFunc: func(old interface{}, new interface{}) {
			glog.Info("UpdateFunc called with objects: %v, %v", old, new)
			key, err := cache.MetaNamespaceKeyFunc(new)
			if err == nil {
				c.queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {
			glog.Info("DeleteFunc called with object: %v", obj)
			// IndexerInformer uses a delta nodeQueue, therefore for deletes we have to use this
			// key function.
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			if err == nil {
				c.queue.Add(key)
			}
		},
	})

	return c
}

// Run will set up the event handlers for types we are interested in, as well
// as syncing informer caches and starting workers. It will block until stopCh
// is closed, at which point it will shutdown the workqueue and wait for
// workers to finish processing their current work items.
func (c *Controller) Run(threadiness int, stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer c.queue.ShutDown()

	// Start the informer factories to begin populating the informer caches
	glog.Info("Starting Foo controller")

	// Wait for the caches to be synced before starting workers
	glog.Info("Waiting for informer caches to sync")
	if ok := cache.WaitForCacheSync(stopCh, c.epSynced, c.deploymentsSynced, c.servicesSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}

	glog.Info("Starting workers")
	// Launch two workers to process the resources
	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, 5*time.Second, stopCh)
	}

	glog.Info("Started workers")
	<-stopCh
	glog.Info("Shutting down workers")

	return nil
}

// runWorker is a long-running function that will continually call the
// processNextWorkItem function in order to read and process a message on the
// workqueue.
func (c *Controller) runWorker() {
	for c.processNextWorkItem() {
	}
}

// processNextWorkItem will read a single work item off the workqueue and
// attempt to process it, by calling the syncHandler.
func (c *Controller) processNextWorkItem() bool {
	obj, shutdown := c.queue.Get()

	if shutdown {
		return false
	}

	// We wrap this block in a func so we can defer c.workqueue.Done.
	err := func(obj interface{}) error {
		// We call Done here so the workqueue knows we have finished
		// processing this item. We also must remember to call Forget if we
		// do not want this work item being re-queued. For example, we do
		// not call Forget if a transient error occurs, instead the item is
		// put back on the workqueue and attempted again after a back-off
		// period.
		defer c.queue.Done(obj)
		var key string
		var ok bool
		// We expect strings to come off the workqueue. These are of the
		// form namespace/name. We do this as the delayed nature of the
		// workqueue means the items in the informer cache may actually be
		// more up to date that when the item was initially put onto the
		// workqueue.
		if key, ok = obj.(string); !ok {
			// As the item in the workqueue is actually invalid, we call
			// Forget here else we'd go into a loop of attempting to
			// process a work item that is invalid.
			c.queue.Forget(obj)
			runtime.HandleError(fmt.Errorf("expected string in workqueue but got %#v", obj))
			return nil
		}
		// Run the syncHandler, passing it the namespace/name string of the
		// Foo resource to be synced.
		if err := c.syncHandler(key); err != nil {
			return fmt.Errorf("error syncing '%s': %s", key, err.Error())
		}
		// Finally, if no error occurs we Forget this item so it does not
		// get queued again until another change happens.
		c.queue.Forget(obj)
		glog.Infof("Successfully synced '%s'", key)
		return nil
	}(obj)

	if err != nil {
		runtime.HandleError(err)
		return true
	}

	return true
}

// syncHandler compares the actual state with the desired, and attempts to
// converge the two
func (c *Controller) syncHandler(key string) error {

	// Convert the namespace/name string into a distinct namespace and name
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}
	fmt.Printf("\nReceived: namespace: %v, name: %v\n", namespace, name)

	ep, err := c.epLister.EventProviders(namespace).Get(name)
	if err != nil {
		return fmt.Errorf("error getting resource: %v", err)
	}
	fmt.Printf("eventprovider: %v", ep)

	switch ep.Spec.ProviderName {
	case "eventgrid.azure.com":
		if ep.Spec.EventType != "Microsoft.Storage" {
			return fmt.Errorf("can only handle storage events")
		}

		// TODO - also check deployment, service and ingress health

		// first check for deployment
		deploymentName := fmt.Sprintf("%s-%s-deployment", ep.Name, ep.Spec.StorageAccount)
		deployment, err := c.deploymentsLister.Deployments(ep.Namespace).Get(deploymentName)
		// If the resource doesn't exist, we'll create it
		if errors.IsNotFound(err) {
			deployment, err = c.kubeclientset.AppsV1().Deployments(ep.Namespace).Create(newDeployment(ep, deploymentName))
		}
		fmt.Printf("deployment name: %v", deployment.Name)

		// If an error occurs during Get/Create, we'll requeue the item so we can
		// attempt processing again later. This could have been caused by a
		// temporary network failure, or any other transient reason.
		if err != nil {
			return err
		}

		// check the service
		serviceName := fmt.Sprintf("%s-%s-service", ep.Name, ep.Spec.StorageAccount)
		service, err := c.servicesLister.Services(ep.Namespace).Get(serviceName)
		// If the resource doesn't exist, we'll create it
		if errors.IsNotFound(err) {
			service, err = c.kubeclientset.CoreV1().Services(ep.Namespace).Create(newService(ep, serviceName, deploymentName))
		}
		fmt.Printf("service name: %v", service.Name)

		// If an error occurs during Get/Create, we'll requeue the item so we can
		// attempt processing again later. This could have been caused by a
		// temporary network failure, or any other transient reason.
		if err != nil {
			return err
		}

		// check ingress
		ingressName := fmt.Sprintf("%s-%s-ingress", ep.Name, ep.Spec.Host)
		ingress, err := c.ingressLister.Ingresses(ep.Namespace).Get(ingressName)
		if errors.IsNotFound(err) {
			ingress, err = c.kubeclientset.ExtensionsV1beta1().Ingresses(ep.Namespace).Create(&v1beta1.Ingress{})
		}
		fmt.Printf("ingress name: %v", ingress.Name)

		// If an error occurs during Get/Create, we'll requeue the item so we can
		// attempt processing again later. This could have been caused by a
		// temporary network failure, or any other transient reason.
		if err != nil {
			return err
		}

		// check eventsubscription exists for given storage account
		exists, err := eventgrid.CheckEventSubscription("", "")
		if err != nil {
			return fmt.Errorf("cannot check eventgrid subscription: %v", err)
		}
		// if the eventsubscription does not exist, create it
		if !exists {
			err = eventgrid.CreateOrUpdateEventSubscription("", "", ep.Spec.StorageAccount, ep.Spec.Host)
			if err != nil {
				return err
			}
		}

	default:
		return fmt.Errorf("cannot handle provider %v", ep.Spec.ProviderName)
	}

	return nil
}

// newDeployment creates a new Deployment based on an eventprovider
func newDeployment(ep *v1alpha1.EventProvider, name string) *appsv1.Deployment {

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: ep.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			// TODO - this is hardcoded
			Replicas: to.Int32Ptr(2),
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  name,
							Image: ep.Spec.HostImage,
						},
					},
				},
			},
		},
	}
}

func newService(ep *v1alpha1.EventProvider, serviceName, deploymentName string) *corev1.Service {
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: serviceName,
		},
		Spec: corev1.ServiceSpec{

			Selector: map[string]string{"app": deploymentName},
			Ports: []corev1.ServicePort{
				{
					Name: "eventgrid-80",
					Port: 80,
					TargetPort: intstr.IntOrString{
						IntVal: 80,
					},
				},
			},
		},
	}
}
